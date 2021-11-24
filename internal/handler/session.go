package handler

import (
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

func (h Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var credentials hmm.LoginCredentials
	err := httpresponse.Unmarshal(r, &credentials)
	if err != nil {
		h.Logger.Errorf("Unmarshal error: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	// create session and cookie
	s, err := h.SessionService.Create(credentials.Email, credentials.Password)
	if err != nil {
		h.Logger.Errorf("unable to create session: %+v", err)
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
