package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dmartzol/hmmm/internal/models"
	"github.com/dmartzol/hmmm/pkg/httpresponse"
	"github.com/dmartzol/hmmm/pkg/randutil"
	"github.com/dmartzol/hmmm/pkg/timeutils"
	"github.com/gorilla/mux"
)

type accountStorage interface {
	Account(id int64) (*models.Account, error)
	Accounts() (models.Accounts, error)
	AccountExists(email string) (bool, error)
	AccountWithCredentials(email, allegedPassword string) (*models.Account, error)
	CreateAccount(first, last, email, password, confirmationCode string, dob time.Time, gender, phone *string) (*models.Account, *models.Confirmation, error)
	CreateConfirmation(accountID int64, t models.ConfirmationType) (*models.Confirmation, error)
	PendingConfirmationByKey(key string) (*models.Confirmation, error)
	FailedConfirmationIncrease(id int64) (*models.Confirmation, error)
	Confirm(id int64) (*models.Confirmation, error)
}

func (api API) GetAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	err := api.AuthorizeAccount(requesterID, models.PermissionAccountsView)
	if err != nil {
		log.Printf("GetAccounts AuthorizeAccount ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accs, err := api.Accounts()
	if err != nil {
		log.Printf("accounts: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, accs.Views(nil))
}

func (api API) GetAccount(w http.ResponseWriter, r *http.Request) {
	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params["id"]
	if !ok {
		http.Error(w, "parameter 'id' not found", http.StatusBadRequest)
		return
	}
	accountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong parameter '%s'", idString), http.StatusBadRequest)
		return
	}

	// checking permissions
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	if requesterID != accountID {
		err := api.AuthorizeAccount(requesterID, models.PermissionAccountsView)
		if err != nil {
			log.Printf("WARNING: account %d requested to see account %d", requesterID, accountID)
			log.Printf("GetAccounts AuthorizeAccount ERROR: %+v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	a, err := api.Account(accountID)
	if err != nil {
		log.Printf("Account: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, a.View(nil))
}

func (api API) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	exists, err := api.AccountExists(req.Email)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		// see: https://stackoverflow.com/questions/9269040/which-http-response-code-for-this-email-is-already-registered
		err = fmt.Errorf("email '%s' already registered", req.Email)
		log.Printf("%+v", err)
		http.Error(w, fmt.Sprintf("email '%s' alrady exists", req.Email), http.StatusConflict)
		return
	}
	err = req.Validate()
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	parsedDOB, err := time.Parse(timeutils.LayoutISO, req.DOB)
	if err != nil {
		log.Printf("%s: %+v", req.DOB, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	code, err := randutil.RandomCode(6)
	if err != nil {
		log.Printf("RandomCode: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	a, _, err := api.storage.CreateAccount(
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("confirmation key: %s", code)

	// create session and cookie
	s, err := api.storage.CreateSession(a.ID)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:   hmmmCookieName,
		Value:  s.SessionIdentifier,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)

	// TODO: send confirmation key by email

	httpresponse.RespondJSON(w, a.View(nil))
}

func (api API) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// TODO: create confirmation key in db
	// TODO: send email with link to reset password
	httpresponse.RespondText(w, "If the account exists, an email will be sent with recovery details.", http.StatusAccepted)
}

func (api API) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// fetching requester id
	requesterID := ctx.Value(contextRequesterAccountIDKey)
	if requesterID == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	a, err := api.Account(requesterID.(int64))
	if err != nil {
		log.Printf("Account: %+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var req models.ConfirmEmailRequest
	err = httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	c, err := api.PendingConfirmationByKey(req.ConfirmationKey)
	if err != nil {
		log.Printf("PendingConfirmationByKey: %+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if c.FailedConfirmationsCount >= 3 {
		log.Printf("FailedConfirmationsCount: %d", c.FailedConfirmationsCount)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// check if user is trying to confirm current email
	if c.ConfirmationTarget == nil {
		log.Printf("confirmation target is null for key %s", req.ConfirmationKey)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if a.Email != *c.ConfirmationTarget {
		log.Printf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// check if keys match
	if c.Key != req.ConfirmationKey {
		_, err := api.FailedConfirmationIncrease(c.ID)
		if err != nil {
			log.Printf("FailedConfirmationIncrease: %+v", err)
		}
		log.Printf("confirmation target %s does not match account email %s", *c.ConfirmationTarget, a.Email)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// confirmation went OK
	_, err = api.Confirm(c.ID)
	if err != nil {
		log.Printf("Confirm: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email has been confirmed.")
}

func (api API) AddAccountRole(w http.ResponseWriter, r *http.Request) {
	var req models.AddAccountRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("AddAccountRole Unmarshal ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params["id"]
	if !ok {
		http.Error(w, "parameter 'id' not found", http.StatusBadRequest)
		return
	}
	requestedAccountID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong parameter '%s'", idString), http.StatusBadRequest)
		return
	}

	role, err := api.Role(req.RoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("AddAccountRole Role ERROR: %+v", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Printf("AddAccountRole Role ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	accRole, err := api.storage.AddAccountRole(role.ID, requestedAccountID)
	if err != nil {
		log.Printf("AddAccountRole storage.AddAccountRole ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, accRole.View(nil))
}
