package api

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/dmartzol/hmm/internal/storage/postgres"
	"github.com/dmartzol/hmm/pkg/httpresponse"
)

func (h API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicRoutes := map[string]string{
			"/v1/version":  "GET",
			"/v1/sessions": "POST",
			"/v1/accounts": "POST",
		}
		method, in := publicRoutes[r.RequestURI]
		if in && method == r.Method {
			next.ServeHTTP(w, r)
			return
		}

		c, err := r.Cookie(hmmmCookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				h.Logger.Info("No cookie found in request")
				httpresponse.RespondJSONError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			default:
				h.Logger.Errorf("error getting cookie: %v", err)
				httpresponse.RespondJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		s, err := h.SessionService.UpdateSession(c.Value)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				h.Logger.Errorf("unable to find session %q: %+v", c.Value, err)
				httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
			case errors.Is(err, postgres.ErrExpiredResource):
				h.Logger.Errorf("session %q is expired: %+v", c.Value, err)
				httpresponse.RespondJSONError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			default:
				log.Printf("error updating session %q: %+v", c.Value, err)
				httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
			}
			return
		}

		// Setting up context
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextRequesterAccountIDKey, s.AccountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
