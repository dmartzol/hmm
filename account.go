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

	"github.com/gorilla/mux"
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var info registerRequest
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	exists, err := db.EmailExists(info.Email)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		// see: https://stackoverflow.com/questions/9269040/which-http-response-code-for-this-email-is-already-registered
		err = fmt.Errorf("email '%s' already registered", info.Email)
		log.Printf("%+v", err)
		http.Error(w, fmt.Sprintf("email '%s' alrady exists", info.Email), http.StatusConflict)
		return
	}
	err = info.validate()
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	parsedDOB, err := time.Parse(layoutISO, info.DOB)
	if err != nil {
		log.Printf("%s: %+v", info.DOB, err)
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
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// create session and cookie
	s, err := db.CreateSession(a.ID)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:   "session",
		Value:  s.SessionIdentifier,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)

	// Create email confirmation code and send email

	json.NewEncoder(w).Encode(s)
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		err := fmt.Errorf("param 'id' not found")
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("%+v", id)
}
