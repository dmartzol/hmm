package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Account represents a user account
type Account struct {
	Row
	FirstName         string `db:"first_name"`
	LastName          string `db:"last_name"`
	DOB               time.Time
	Gender            *string
	Active            bool
	FailedLoginsCount int64   `db:"failed_logins_count"`
	DoorCode          *string `db:"door_code"`
	PassHash          string
	Email             string
	PhoneNumber       *string `db:"phone_number"`

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

func createAccount(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	exists, err := db.EmailExists(req.Email)
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
	err = req.validate()
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	parsedDOB, err := time.Parse(layoutISO, req.DOB)
	if err != nil {
		log.Printf("%s: %+v", req.DOB, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	a, err := db.CreateAccount(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
		parsedDOB,
		req.Gender,
		req.PhoneNumber,
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
		Name:   hackerSpaceCookieName,
		Value:  s.SessionIdentifier,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)

	// TODO: Create email confirmation code and send it

	json.NewEncoder(w).Encode(s)
}

func getAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idString, ok := params["id"]
	if !ok {
		http.Error(w, "parameter 'id' not found", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong parameter '%s'", idString), http.StatusBadRequest)
		return
	}
	a, err := db.Account(id)
	if err != nil {
		log.Printf("Account: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a)
}

type resetPasswordRequest struct {
	Email string
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	err := Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// TODO: create confirmation code in db
	// TODO: send email with link to reset password
	HTTPRespond(w, "If the account exists, an email will be sent with recovery details.", http.StatusAccepted)
}

type confirmEmailRequest struct {
	ConfirmationCode string
}

func confirmEmail(w http.ResponseWriter, r *http.Request) {
	var req confirmEmailRequest
	err := Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email has been confirmed.")
}
