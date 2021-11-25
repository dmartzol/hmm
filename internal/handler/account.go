package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dmartzol/hmm/internal/api"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/pkg/httpresponse"
	"github.com/dmartzol/hmm/pkg/randutil"
	"github.com/gorilla/mux"
)

func (h Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
	if err != nil {
		h.Logger.Errorf("unable to authorize account %d: %+v", requesterID, err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	accs, err := h.AccountService.Accounts()
	if err != nil {
		h.Logger.Errorf("unable to fetch accounts: %v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	h.AccountService.PopulateAccounts(accs)

	httpresponse.RespondJSON(w, api.AccountsView(accs, nil))
}

func (h Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		h.Logger.Errorf("unable to parse %q: %v", idString, err)
		httpresponse.RespondJSONError(w, fmt.Sprintf("wrong parameter '%s'", idString), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	if requesterID != accountID {
		err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
		if err != nil {
			h.Logger.Errorf("account %d requested to see account %d: %v", requesterID, accountID, err)
			httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
			return
		}
	}

	a, err := h.AccountService.Account(accountID)
	if err == sql.ErrNoRows {
		log.Printf("account %d not found", accountID)
		httpresponse.RespondJSONError(w, "", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("could not fetch account %d: %+v", accountID, err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	h.AccountService.PopulateAccount(a)

	httpresponse.RespondJSON(w, api.AccountView(a, nil))
}

func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req api.RegisterRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	err = req.ValidateAndNormalize()
	if err != nil {
		h.Logger.Errorf("error validating: %+v", req.Email)
		httpresponse.RespondJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, err := randutil.RandomCode(6)
	if err != nil {
		h.Logger.Errorf("error generating random code for %q: %+v", req.Email, err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	inputAccount := hmm.Account{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DOB:         req.DOBTime,
		PhoneNumber: req.PhoneNumber,
	}
	a, _, err := h.AccountService.Create(&inputAccount, req.Password, code)
	if err != nil {
		// TODO: respond with 409 on existing email address
		// see: https://stackoverflow.com/questions/9269040/which-http-response-code-for-this-email-is-already-registered
		h.Logger.Errorf("error creating account: %+v", req.Email)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("confirmation key: %s", code)

	// create session and cookie
	s, err := h.SessionService.Create(a.Email, req.Password)
	if err != nil {
		h.Logger.Errorf("error creating session: %+v", req.Email)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:   hmmmCookieName,
		Value:  s.Token,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)

	// TODO: send confirmation key by email

	httpresponse.RespondJSON(w, api.AccountView(a, nil))
}

func (h Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req hmm.ResetPasswordRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	// TODO: create confirmation key in db
	// TODO: send email with link to reset password
	httpresponse.RespondText(w, "not implemented", http.StatusNotImplemented)
}

func (h Handler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey)
	if requesterID == nil {
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	a, err := h.AccountService.Account(requesterID.(int64))
	if err != nil {
		log.Printf("Account: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	var req hmm.ConfirmEmailRequest
	err = httpresponse.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	c, err := h.ConfirmationService.PendingConfirmationByKey(req.ConfirmationKey)
	if err != nil {
		h.Logger.Errorf("failed to fetch confirmation by key: %v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	if c.FailedConfirmationsCount >= 3 {
		h.Logger.Errorf("too many attempts to confirm", err)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	// check if user is trying to confirm current email
	if c.ConfirmationTarget == nil {
		h.Logger.Errorf("confirmation target is null for key %s", req.ConfirmationKey)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	if a.Email != *c.ConfirmationTarget {
		h.Logger.Errorf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	// check if keys match
	if c.Key != req.ConfirmationKey {
		_, err := h.ConfirmationService.FailedConfirmationIncrease(c.ID)
		if err != nil {
			h.Logger.Errorf("failed confirmation increase: %v", err)
		}
		h.Logger.Errorf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	_, err = h.ConfirmationService.Confirm(c.ID)
	if err != nil {
		h.Logger.Errorf("failed to confirm: %v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email has been confirmed.")
}

func (h Handler) AddAccountRole(w http.ResponseWriter, r *http.Request) {
	var req hmm.AddAccountRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("AddAccountRole Unmarshal ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMSg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		httpresponse.RespondJSONError(w, errMSg, http.StatusInternalServerError)
		return
	}
	requestedAccountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}

	role, err := h.RoleService.Role(req.RoleID)
	if err != nil {
		log.Printf("AddAccountRole ERROR fetching role: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	accRole, err := h.RoleService.AddRoleToAccount(role.ID, requestedAccountID)
	if err != nil {
		log.Printf("AddAccountRole storage.AddAccountRole ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, api.AccountRoleView(accRole, nil))
}

func (h Handler) GetAccountRoles(w http.ResponseWriter, r *http.Request) {
	var req hmm.AddAccountRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("GetAccountRoles Unmarshal ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMsg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}
	requestedAccountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}

	// checking permissions
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	if requesterID != requestedAccountID {
		err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
		if err != nil {
			log.Printf("WARNING: account %d requested to see account %d", requesterID, requestedAccountID)
			log.Printf("GetAccounts AuthorizeAccount ERROR: %+v", err)
			httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
			return
		}
	}

	rs, err := h.RoleService.RolesForAccount(requestedAccountID)
	if err != nil {
		log.Printf("GetAccountRoles RolesForAccount ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, api.RolesView(rs, nil))
}
