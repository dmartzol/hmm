package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dmartzol/hackerspace/internal/models"
	"github.com/dmartzol/hackerspace/pkg/httpresponse"
	"github.com/dmartzol/hackerspace/pkg/timeutils"
	"github.com/gorilla/mux"
)

func createAccount(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	err := httpresponse.Unmarshal(r, &req)
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

func resetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// TODO: create confirmation code in db
	// TODO: send email with link to reset password
	httpresponse.Respond(w, "If the account exists, an email will be sent with recovery details.", http.StatusAccepted)
}

func confirmEmail(w http.ResponseWriter, r *http.Request) {
	var req models.ConfirmEmailRequest
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("JSON: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email has been confirmed.")
}
