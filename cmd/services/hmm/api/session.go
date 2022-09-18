package api

import (
	"errors"
	"net/http"

	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/internal/storage/postgres"
	"github.com/dmartzol/hmm/pkg/httpresponse"
)

const (
	// sessionLength represents the duration(in seconds) a session will be valid for
	sessionLength = 345600
)

func (h API) GetSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hmmmCookieName)
	if err != nil {
		h.Logger.Errorf("unable to fetch cookie: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	s, err := h.SessionService.SessionFromToken(c.Value)
	if err != nil {
		h.Logger.Errorf("unable to fetch session: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	httpresponse.RespondJSON(w, s.View(nil))
}

func (h API) CreateSession(w http.ResponseWriter, r *http.Request) {
	var credentials hmm.LoginCredentials
	err := httpresponse.Unmarshal(r, &credentials)
	if err != nil {
		h.Logger.Errorf("Unmarshal error: %+v", err)
		httpresponse.RespondJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// create session
	s, err := h.SessionService.Create(credentials.Email, credentials.Password)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrInvalidCredentials):
			h.Logger.Warn("invalid credentials")
			httpresponse.RespondJSONError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		default:
			h.Logger.Errorf("unable to create session: %+v", err)
			httpresponse.RespondJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// creating and setting cookie
	cookie := &http.Cookie{
		Name:   hmmmCookieName,
		Value:  s.Token,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)
	httpresponse.RespondJSON(w, s.View(nil))
}

func (h API) ExpireSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hmmmCookieName)
	if err != nil {
		h.Logger.Errorf("error fetching cookie: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	session, err := h.SessionService.ExpireSession(c.Value)
	if err != nil {
		h.Logger.Errorf("unable to expire session: %+v", err)
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
