package api

import (
	"errors"
	"net/http"

	"github.com/dmartzol/hmm/internal/dao/postgres"
	"github.com/dmartzol/hmm/internal/hmm"
)

const (
	// sessionLength represents the duration(in seconds) a session will be valid for
	sessionLength = 345600
	hmmCookieName = "Hmm-Cookie"
)

func (h API) GetSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hmmCookieName)
	if err != nil {
		h.Logger.Errorf("unable to fetch cookie: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	s, err := h.SessionService.SessionFromToken(c.Value)
	if err != nil {
		h.Logger.Errorf("unable to fetch session: %+v", err)
		h.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}
	h.RespondJSON(w, s.View(nil))
}

func (h API) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var credentials hmm.LoginCredentials
	err := h.Unmarshal(r, &credentials)
	if err != nil {
		h.Logger.Errorf("Unmarshal error: %+v", err)
		h.RespondJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// create session
	s, err := h.SessionService.Create(ctx, credentials.Email, credentials.Password)
	if err != nil {
		switch {
		case errors.Is(err, postgres.ErrInvalidCredentials):
			h.Logger.Warn("invalid credentials")
			h.RespondJSONError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		default:
			h.Logger.Errorf("unable to create session: %+v", err)
			h.RespondJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// creating and setting cookie
	cookie := &http.Cookie{
		Name:   hmmCookieName,
		Value:  s.Token,
		MaxAge: sessionLength,
	}
	http.SetCookie(w, cookie)
	h.RespondJSON(w, s.View(nil))
}

func (h API) ExpireSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(hmmCookieName)
	if err != nil {
		h.Logger.Errorf("error fetching cookie: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	session, err := h.SessionService.ExpireSession(c.Value)
	if err != nil {
		h.Logger.Errorf("unable to expire session: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	c = &http.Cookie{
		Name:   hmmCookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	h.RespondJSON(w, session.View(nil))
}
