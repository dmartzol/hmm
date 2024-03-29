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
	"github.com/dmartzol/hmm/pkg/timeutils"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type CreateAccountRequest struct {
	FirstName   string
	LastName    string
	DOBString   string    `json:"dob"`
	DOB         time.Time `json:"-"`
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
			&c.DOBString,
			validation.Required,
			validation.Date(time.RFC3339),
		),
		validation.Field(
			&c.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&c.Password,
			validation.Required,
			validation.Length(10, 500),
		),
	)
}

func (r *CreateAccountRequest) normalize() error {
	r.FirstName = normalizeName(r.FirstName)
	r.LastName = normalizeName(r.LastName)
	var err error
	r.DOB, err = time.Parse(time.RFC3339, r.DOBString)
	if err != nil {
		return fmt.Errorf("error parsing DOB %q: %w", r.DOBString, err)
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
	ctx := r.Context()
	ctx, span := otel.Tracer(appName).Start(ctx, "Resources.CreateAccount")
	defer span.End()

	var req CreateAccountRequest
	err := re.Unmarshal(r, &req)
	if err != nil {
		re.Logger.Errorf("unable to unmarshal: %+v", err)
		re.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	err = req.ValidateAndNormalize()
	if err != nil {
		re.Logger.Errorf("error validating payload: %+v", err)
		re.RespondJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// generate a random confirmation code and password
	randomConfirmationCode := RandomConfirmationCode(6)

	// we use a hmm.Account here because the db library does not have access to the CreateAccountRequest type
	inputAccount := hmm.Account{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DOB:         req.DOB,
		PhoneNumber: req.PhoneNumber,
	}
	a, _, err := re.AccountService.Create(ctx, &inputAccount, req.Password, randomConfirmationCode)
	if err != nil {
		// TODO: respond with 409 on existing email address
		// see: https://stackoverflow.com/questions/9269040/which-http-response-code-for-this-email-is-already-registered
		re.Logger.Errorf("error creating account: %+v", err)
		re.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	// create a session for the new account
	s, err := re.SessionService.Create(ctx, a.Email, req.Password)
	if err != nil {
		re.Logger.Errorf("error creating session: %+v", req.Email)
		re.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:   hmmCookieName,
		Value:  s.Token,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)

	// TODO: send confirmation key by email

	re.RespondJSON(w, AccountView(a, nil))
}

func (h API) GetAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
	if err != nil {
		h.Logger.Errorf("unable to authorize account %d: %+v", requesterID, err)
		h.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	accs, err := h.AccountService.Accounts()
	if err != nil {
		h.Logger.Errorf("unable to fetch accounts: %v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	h.AccountService.PopulateAccounts(accs)

	h.RespondJSON(w, AccountsView(accs, nil))
}

func (h API) GetAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		h.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		h.Logger.Errorf("unable to parse %q: %v", idString, err)
		h.RespondJSONError(w, fmt.Sprintf("wrong parameter '%s'", idString), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	if requesterID != accountID {
		err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
		if err != nil {
			h.Logger.Errorf("account %d requested to see account %d: %v", requesterID, accountID, err)
			h.RespondJSONError(w, "", http.StatusUnauthorized)
			return
		}
	}

	a, err := h.AccountService.Account(accountID)
	if err == sql.ErrNoRows {
		log.Printf("account %d not found", accountID)
		h.RespondJSONError(w, "", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("could not fetch account %d: %+v", accountID, err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	h.AccountService.PopulateAccount(a)

	h.RespondJSON(w, AccountView(a, nil))
}

func (h API) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req hmm.ResetPasswordRequest
	err := h.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	// TODO: create confirmation key in db
	// TODO: send email with link to reset password
	h.RespondText(w, "not implemented", http.StatusNotImplemented)
}

func (h API) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey)
	if requesterID == nil {
		h.Logger.Errorf("no requester ID in context")
		h.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	a, err := h.AccountService.Account(requesterID.(int64))
	if err != nil {
		h.Logger.Errorf("unable to fetch account %d: %+v", requesterID.(int64), err)
		h.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	var req hmm.ConfirmEmailRequest
	err = h.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	c, err := h.ConfirmationService.PendingConfirmationByKey(req.ConfirmationKey)
	if err != nil {
		h.Logger.Errorf("failed to fetch confirmation by key: %v", err)
		h.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	if c.FailedConfirmationsCount >= 3 {
		h.Logger.Errorf("too many attempts to confirm", err)
		h.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

	if c.ConfirmationTarget == nil {
		h.Logger.Errorf("confirmation target is null for key %s", req.ConfirmationKey)
		h.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	// check if user is trying to confirm current email
	if a.Email != *c.ConfirmationTarget {
		h.Logger.Errorf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		h.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

	if c.Key != req.ConfirmationKey {
		_, err := h.ConfirmationService.FailedConfirmationIncrease(c.ID)
		if err != nil {
			h.Logger.Errorf("failed confirmation increase: %v", err)
		}
		h.Logger.Errorf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		h.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

	_, err = h.ConfirmationService.Confirm(c.ID)
	if err != nil {
		h.Logger.Errorf("failed to confirm: %v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h API) AddAccountRole(w http.ResponseWriter, r *http.Request) {
	var req hmm.AddAccountRoleReq
	err := h.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMSg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		h.Logger.Errorf(errMSg)
		h.RespondJSONError(w, errMSg, http.StatusInternalServerError)
		return
	}

	requestedAccountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		h.Logger.Errorf(errMsg)
		h.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}

	role, err := h.RoleService.Role(req.RoleID)
	if err != nil {
		h.Logger.Errorf("unable to fetch role %d: %+v", req.RoleID, err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	accRole, err := h.RoleService.AddRoleToAccount(role.ID, requestedAccountID)
	if err != nil {
		h.Logger.Errorf("unable to add role %d to account %d: %+v", role.ID, requestedAccountID, err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	h.RespondJSON(w, AccountRoleView(accRole, nil))
}

func (h API) GetAccountRoles(w http.ResponseWriter, r *http.Request) {
	var req hmm.AddAccountRoleReq
	err := h.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMsg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		h.Logger.Errorf(errMsg)
		h.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}
	requestedAccountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		h.Logger.Errorf(errMsg)
		h.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	if requesterID != requestedAccountID {
		err := h.AuthorizeAccount(requesterID, hmm.PermissionAccountsView)
		if err != nil {
			h.Logger.Errorf("unable to authorize account %d: %+v", requesterID, err)
			h.RespondJSONError(w, "", http.StatusUnauthorized)
			return
		}
	}

	rs, err := h.RoleService.RolesForAccount(requestedAccountID)
	if err != nil {
		h.Logger.Errorf("unable to fetch roles for account %d: %+v", requestedAccountID, err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	h.RespondJSON(w, RolesView(rs, nil))
}
