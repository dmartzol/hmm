package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmartzol/hackerspace/internal/models"
	"github.com/dmartzol/hackerspace/pkg/httpresponse"
)

const (
	// sessionLength represents the duration(in minutes) a session will be valid for
	sessionLength = 3600
)

func createSession(w http.ResponseWriter, r *http.Request) {
	var credentials models.LoginCredentials
	err := httpresponse.Unmarshal(r, &credentials)
	if err != nil {
		log.Printf("Unmarshal: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fetching account with credentials(errors reurned should be purposedly broad)
	registered, err := db.EmailExists(credentials.Email)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !registered {
		log.Printf("unable to find email '%s' in db", credentials.Email)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	a, err := db.AccountWithCredentials(credentials.Email, credentials.Password)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	credentials.Password = ""

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
	json.NewEncoder(w).Encode(s)
}

func deleteSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hackerSpaceCookieName)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = db.DeleteSession(c.Value)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	c = &http.Cookie{
		Name:   hackerSpaceCookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	httpresponse.Respond(w, "Session deleted.", http.StatusOK)
}
