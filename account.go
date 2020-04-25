package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Account represents a user account
type Account struct {
	Row
	FirstName           string `db:"first_name"`
	LastName            string `db:"last_name"`
	DOB                 time.Time
	Gender              *string
	Active              bool
	FailedLoginAttempts int64   `db:"failed_login_attempts"`
	DoorCode            *string `db:"door_code"`
	PassHash            string
	Email               string
	PhoneNumber         *string `db:"phone_number"`

	RoleID                    *int64 `db:"role_id"`
	EmailID                   int64  `db:"email_id"`
	PhoneNumberID             *int64 `db:"phone_number_id"`
	ExternalPaymentCustomerID *int64 `db:"external_payment_customer_id"`
}

type registerRequest struct {
	FirstName   string
	LastName    string
	DOB         string
	Gender      *string
	PhoneNumber *string
	Email       string
	Password    string
}

func (r registerRequest) validate() error {
	if r.FirstName == "" {
		return errors.New("must provide first name")
	}
	if r.LastName == "" {
		return errors.New("must provide last name")
	}
	if len(r.Password) < 6 {
		return errors.New("password too short")
	}
	_, err := time.Parse(layoutISO, r.DOB)
	if err != nil {
		return fmt.Errorf("time.Parse %v: %w", r.DOB, err)
	}
	return nil
}

func validEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	if !strings.Contains(email, ".") {
		return false
	}
	return true
}

func register(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var info registerRequest
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("Unmarshal: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = info.validate()
	if err != nil {
		log.Printf("validate: %+v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	parsedDOB, err := time.Parse(layoutISO, info.DOB)
	if err != nil {
		log.Printf("time.Parse %s: %+v", info.DOB, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	a, err := db.CreateAccount(
		info.FirstName,
		info.LastName,
		info.Email,
		info.Password,
		parsedDOB,
		info.Gender,
		info.PhoneNumber,
	)
	if err != nil {
		log.Printf("CreateAccount: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Create email confirmation code and send email
	json.NewEncoder(w).Encode(a)
}

type loginCredentials struct {
	Email    string
	Password string
}

func login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ReadAll: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var credentials loginCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		log.Printf("Unmarshal: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fetching account with credentials(errors reurned should be purposedly broad)
	registered, err := db.EmailExists(credentials.Email)
	if err != nil {
		log.Printf("EmailExists: %+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if !registered {
		log.Printf("unable to find email '%s' in db", credentials.Email)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	a, err := db.AccountWithCredentials(credentials.Email, credentials.Password)
	if err != nil {
		log.Printf("AccountWithCredentials: %+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	credentials.Password = ""

	// create session and cookie
	s, err := db.CreateSession(a.ID, time.Now().Add(time.Minute*5), uuid.New().String())
	if err != nil {
		log.Printf("CreateSession: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Printf("Cookie: %+v", err)
		cookie = &http.Cookie{
			Name:  "session",
			Value: s.Token,
		}
		http.SetCookie(w, cookie)
	}
	json.NewEncoder(w).Encode(s)
}

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	s, err := db.SessionFromToken(c.Value)
	if err != nil {
		return false
	}
	if s.ExpirationDate.Before(time.Now()) {
		return false
	}
	return true
}

func logout(w http.ResponseWriter, r *http.Request) {

}
