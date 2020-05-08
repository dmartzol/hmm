package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dmartzol/hackerspace/internal/models"
	"github.com/dmartzol/hackerspace/pkg/httpresponse"
)

const (
	// sessionLength represents the duration(in minutes) a session will be valid for
	sessionLength = 3600
)

type sessionStorage interface {
	SessionFromIdentifier(identifier string) (*models.Session, error)
	CreateSession(accountID int64) (*models.Session, error)
	DeleteSession(identifier string) error
	CleanSessionsOlderThan(age time.Duration) (int64, error)
	UpdateSession(identifier string) (*models.Session, error)
}

func (api API) createSession(w http.ResponseWriter, r *http.Request) {
	var credentials models.LoginCredentials
	err := httpresponse.Unmarshal(r, &credentials)
	if err != nil {
		log.Printf("Unmarshal: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fetching account with credentials(errors reurned should be purposedly broad)
	registered, err := api.AccountExists(credentials.Email)
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
	a, err := api.AccountWithCredentials(credentials.Email, credentials.Password)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	credentials.Password = ""

	// create session and cookie
	s, err := api.CreateSession(a.ID)
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

func (api API) deleteSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hackerSpaceCookieName)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = api.DeleteSession(c.Value)
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
