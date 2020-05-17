package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/dmartzol/hmmm/internal/models"
	"github.com/dmartzol/hmmm/pkg/httpresponse"
)

const (
	// sessionLength represents the duration(in seconds) a session will be valid for
	sessionLength = 3600
)

type sessionStorage interface {
	SessionFromIdentifier(identifier string) (*models.Session, error)
	CreateSession(accountID int64) (*models.Session, error)
	DeleteSession(identifier string) error
	CleanSessionsOlderThan(age time.Duration) (int64, error)
	// UpdateSession updates a session in the db with the current timestamp
	UpdateSession(identifier string) (*models.Session, error)
}

func (api API) CreateSession(w http.ResponseWriter, r *http.Request) {
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
	s, err := api.storage.CreateSession(a.ID)
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
	httpresponse.RespondText(w, "Successfully logged in", http.StatusOK)
}

func (api API) DeleteSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hackerSpaceCookieName)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = api.storage.DeleteSession(c.Value)
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
	httpresponse.RespondText(w, "Session deleted.", http.StatusOK)
}
