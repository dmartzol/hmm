package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dmartzol/hmm/internal/api"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/pkg/httpresponse"
	"github.com/dmartzol/hmm/pkg/randutil"
	"github.com/dmartzol/hmm/pkg/timeutils"
	"github.com/gorilla/mux"
)

func (h Handler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
	if err != nil {
		log.Printf("GetAccounts AuthorizeAccount ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	accs, err := h.db.Accounts()
	if err != nil {
		log.Printf("accounts: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	h.db.PopulateAccounts(accs)
	httpresponse.RespondJSON(w, api.AccountsToAPI(accs, nil))
}

func (h Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	accountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		httpresponse.RespondJSONError(w, fmt.Sprintf("wrong parameter '%s'", idString), http.StatusBadRequest)
		return
	}

	// checking permissions
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	if requesterID != accountID {
		err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
		if err != nil {
			log.Printf("WARNING: account %d requested to see account %d", requesterID, accountID)
			log.Printf("GetAccounts AuthorizeAccount ERROR: %+v", err)
			httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
			return
		}
	}

	a, err := h.AccountService.Account(accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("account %d not found", accountID)
			httpresponse.RespondJSONError(w, "", http.StatusNotFound)
			return
		} else {
			log.Printf("could not fetch account %d: %+v", accountID, err)
			httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
			return
		}
	}
	h.db.PopulateAccount(a)
	httpresponse.RespondJSON(w, api.AccountToAPI(a, nil))
}

func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req api.RegisterRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	exists, err := h.db.AccountExists(req.Email)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	if exists {
		// see: https://stackoverflow.com/questions/9269040/which-http-response-code-for-this-email-is-already-registered
		err = fmt.Errorf("email '%s' already registered", req.Email)
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, fmt.Sprintf("account with email '%s' alrady exists", req.Email), http.StatusBadRequest)
		return
	}
	// normalizing gender
	if req.Gender != nil {
		if *req.Gender == "female" {
			*req.Gender = "F"
		}
		if *req.Gender == "male" {
			*req.Gender = "M"
		}
	}
	err = req.Validate()
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	parsedDOB, err := time.Parse(timeutils.LayoutISODay, req.DOB)
	if err != nil {
		log.Printf("%s: %+v", req.DOB, err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	code, err := randutil.RandomCode(6)
	if err != nil {
		log.Printf("RandomCode: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	a, _, err := h.db.CreateAccount(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		code,
		parsedDOB,
		req.Gender,
		req.PhoneNumber,
	)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	log.Printf("confirmation key: %s", code)

	// create session and cookie
	s, err := h.db.CreateSession(a.ID)
	if err != nil {
		log.Printf("%+v", err)
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

	httpresponse.RespondJSON(w, api.AccountToAPI(a, nil))
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
	httpresponse.RespondText(w, "If the account exists, an email will be sent with recovery details.", http.StatusAccepted)
}

func (h Handler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// fetching requester id
	requesterID := ctx.Value(contextRequesterAccountIDKey)
	if requesterID == nil {
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	a, err := h.db.Account(requesterID.(int64))
	if err != nil {
		log.Printf("Account: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	var req hmm.ConfirmEmailRequest
	err = httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	c, err := h.db.PendingConfirmationByKey(req.ConfirmationKey)
	if err != nil {
		log.Printf("PendingConfirmationByKey: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	if c.FailedConfirmationsCount >= 3 {
		log.Printf("FailedConfirmationsCount: %d", c.FailedConfirmationsCount)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	// check if user is trying to confirm current email
	if c.ConfirmationTarget == nil {
		log.Printf("confirmation target is null for key %s", req.ConfirmationKey)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	if a.Email != *c.ConfirmationTarget {
		log.Printf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	// check if keys match
	if c.Key != req.ConfirmationKey {
		_, err := h.db.FailedConfirmationIncrease(c.ID)
		if err != nil {
			log.Printf("FailedConfirmationIncrease: %+v", err)
		}
		log.Printf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	// confirmation went OK
	_, err = h.db.Confirm(c.ID)
	if err != nil {
		log.Printf("Confirm: %+v", err)
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

	role, err := h.db.Role(req.RoleID)
	if err != nil {
		log.Printf("AddAccountRole ERROR fetching role: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	accRole, err := h.db.AddAccountRole(role.ID, requestedAccountID)
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

	rs, err := h.db.RolesForAccount(requestedAccountID)
	if err != nil {
		log.Printf("GetAccountRoles RolesForAccount ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, api.RolesView(rs, nil))
}
