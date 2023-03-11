package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/pkg/httpresponse"
	"github.com/dmartzol/hmm/pkg/randutil"
	"github.com/dmartzol/hmm/pkg/timeutils"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
)

type CreateAccountRequest struct {
	FirstName   string
	LastName    string
	DOB         string
	DOBTime     time.Time
	Gender      *string
	PhoneNumber *string
	Email       string
	Password    string
}

// Account is the restricted response body of hmm.Account
// see: https://stackoverflow.com/questions/46427723/golang-elegant-way-to-omit-a-json-property-from-being-serialized
type Account struct {
	ID                         int64 `json:"ID"`
	FirstName, LastName, Email string
	DOB                        string `json:"DateOfBird"`
	PhoneNumber                string
	DoorCode                   string
	Gender                     string
	Active                     bool
	ConfirmedEmail             bool
	ConfirmedPhone             bool
	FailedLoginsCount          int64
	Roles                      []Role
}

func AccountView(a *hmm.Account, options map[string]bool) Account {
	view := Account{
		ID:                a.ID,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		DOB:               a.DOB.Format(timeutils.LayoutISODay),
		Active:            a.Active,
		FailedLoginsCount: a.FailedLoginsCount,
		Email:             a.Email,
		ConfirmedEmail:    a.ConfirmedEmail,
		ConfirmedPhone:    a.ConfirmedPhone,
	}
	if a.DoorCode != nil && options["door_code"] {
		view.DoorCode = *a.DoorCode
	}
	if a.PhoneNumber != nil && options["phone_number"] {
		view.PhoneNumber = *a.PhoneNumber
	}
	if a.Gender != nil {
		view.Gender = *a.Gender
	}
	if a.Roles != nil {
		for _, r := range a.Roles {
			view.Roles = append(view.Roles, RoleView(r, nil))
		}
	}
	return view
}

func AccountsView(accs hmm.Accounts, options map[string]bool) []Account {
	var l []Account
	for _, a := range accs {
		l = append(l, AccountView(a, options))
	}
	return l
}

func (r *CreateAccountRequest) ValidateAndNormalize() error {
	err := r.validate()
	if err != nil {
		return fmt.Errorf("error validating: %w", err)
	}

	err = r.normalize()
	if err != nil {
		return fmt.Errorf("error normalizing: %w", err)
	}

	return nil
}

func (c *CreateAccountRequest) validate() error {
	return validation.ValidateStruct(c,
		validation.Field(
			&c.FirstName,
			validation.Required,
		),
		validation.Field(
			&c.LastName,
			validation.Required,
		),
		validation.Field(
			&c.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&c.Password,
			validation.Required,
			validation.Length(9, 500),
		),
	)
}

func (r *CreateAccountRequest) normalize() error {
	r.FirstName = NormalizeName(r.FirstName)
	r.LastName = NormalizeName(r.LastName)
	var err error
	r.DOBTime, err = time.Parse(timeutils.LayoutISODay, r.DOB)
	if err != nil {
		return fmt.Errorf("error parsing DOB %q: %w", r.DOB, err)
	}
	// normalizing gender only if it is provided
	if r.Gender != nil {
		if strings.EqualFold(*r.Gender, "female") {
			*r.Gender = "F"
		}
		if strings.EqualFold(*r.Gender, "male") {
			*r.Gender = "M"
		}
	}
	return nil
}

func (re Resources) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		re.Logger.Errorf("unable to unmarshal: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	err = req.ValidateAndNormalize()
	if err != nil {
		re.Logger.Errorf("error validating: %+v", req.Email)
		httpresponse.RespondJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, err := randutil.RandomCode(6)
	if err != nil {
		re.Logger.Errorf("error generating random code for %q: %+v", req.Email, err)
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
	a, _, err := re.AccountService.Create(&inputAccount, req.Password, code)
	if err != nil {
		// TODO: respond with 409 on existing email address
		// see: https://stackoverflow.com/questions/9269040/which-http-response-code-for-this-email-is-already-registered
		re.Logger.Errorf("error creating account: %+v", req.Email)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	re.Logger.Infof("confirmation key: %s", code)

	s, err := re.SessionService.Create(a.Email, req.Password)
	if err != nil {
		re.Logger.Errorf("error creating session: %+v", req.Email)
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

	httpresponse.RespondJSON(w, AccountView(a, nil))
}

func (h API) GetAccounts(w http.ResponseWriter, r *http.Request) {
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

	httpresponse.RespondJSON(w, AccountsView(accs, nil))
}

func (h API) GetAccount(w http.ResponseWriter, r *http.Request) {
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

	httpresponse.RespondJSON(w, AccountView(a, nil))
}

func (h API) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req hmm.ResetPasswordRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	// TODO: create confirmation key in db
	// TODO: send email with link to reset password
	httpresponse.RespondText(w, "not implemented", http.StatusNotImplemented)
}

func (h API) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey)
	if requesterID == nil {
		h.Logger.Errorf("no requester ID in context")
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	a, err := h.AccountService.Account(requesterID.(int64))
	if err != nil {
		h.Logger.Errorf("unable to fetch account %d: %+v", requesterID.(int64), err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	var req hmm.ConfirmEmailRequest
	err = httpresponse.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
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

	if c.ConfirmationTarget == nil {
		h.Logger.Errorf("confirmation target is null for key %s", req.ConfirmationKey)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	// check if user is trying to confirm current email
	if a.Email != *c.ConfirmationTarget {
		h.Logger.Errorf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

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
}

func (h API) AddAccountRole(w http.ResponseWriter, r *http.Request) {
	var req hmm.AddAccountRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMSg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		h.Logger.Errorf(errMSg)
		httpresponse.RespondJSONError(w, errMSg, http.StatusInternalServerError)
		return
	}

	requestedAccountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		h.Logger.Errorf(errMsg)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}

	role, err := h.RoleService.Role(req.RoleID)
	if err != nil {
		h.Logger.Errorf("unable to fetch role %d: %+v", req.RoleID, err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	accRole, err := h.RoleService.AddRoleToAccount(role.ID, requestedAccountID)
	if err != nil {
		h.Logger.Errorf("unable to add role %d to account %d: %+v", role.ID, requestedAccountID, err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	httpresponse.RespondJSON(w, AccountRoleView(accRole, nil))
}

func (h API) GetAccountRoles(w http.ResponseWriter, r *http.Request) {
	var req hmm.AddAccountRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMsg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		h.Logger.Errorf(errMsg)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}
	requestedAccountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		h.Logger.Errorf(errMsg)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	if requesterID != requestedAccountID {
		err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
		if err != nil {
			h.Logger.Errorf("unable to authorize account %d: %+v", requesterID, err)
			httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
			return
		}
	}

	rs, err := h.RoleService.RolesForAccount(requestedAccountID)
	if err != nil {
		h.Logger.Errorf("unable to fetch roles for account %d: %+v", requestedAccountID, err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	httpresponse.RespondJSON(w, RolesView(rs, nil))
}
