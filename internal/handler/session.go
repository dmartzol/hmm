package handler

import (
	"log"
	"net/http"

	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/pkg/httpresponse"
)

const (
	// sessionLength represents the duration(in seconds) a session will be valid for
	sessionLength = 345600
)

func (h Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hmmmCookieName)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	s, err := h.db.SessionFromToken(c.Value)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	httpresponse.RespondJSON(w, s.View(nil))
}

func (h Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var credentials hmm.LoginCredentials
	err := httpresponse.Unmarshal(r, &credentials)
	if err != nil {
		log.Printf("Unmarshal: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	// fetching account with credentials(errors reurned should be purposedly broad)
	registered, err := h.db.AccountExists(credentials.Email)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	if !registered {
		log.Printf("unable to find email '%s' in db", credentials.Email)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	a, err := h.db.AccountWithCredentials(credentials.Email, credentials.Password)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	credentials.Password = ""

	// create session and cookie
	s, err := h.db.CreateSession(a.ID)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:   hmmmCookieName,
		Value:  s.Token,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)
	httpresponse.RespondJSON(w, s.View(nil))
}

func (h Handler) ExpireSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hmmmCookieName)
	if err != nil {
		log.Printf("%+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	session, err := h.db.ExpireSessionFromToken(c.Value)
	if err != nil {
		log.Printf("ExpireSession - ERROR expiring session: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	c = &http.Cookie{
		Name:   hmmmCookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	httpresponse.RespondJSON(w, session.View(nil))
}
