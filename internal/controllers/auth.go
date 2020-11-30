package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/dmartzol/hmm/internal/storage/postgres"
	"github.com/dmartzol/hmm/pkg/httpresponse"
)

func (api API) AuthMiddleware(next http.Handler) http.Handler {
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
			log.Printf("AuthMiddleware ERROR getting cookie: %+v", err)
			httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
			return
		}
		s, err := api.db.UpdateSession(c.Value)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("AuthMiddleware ERROR unable to find session %s: %+v", c.Value, err)
				httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
				return
			}
			if err != postgres.ErrExpiredResource {
				log.Printf("AuthMiddleware ERROR session %s is expired: %+v", c.Value, err)
				httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
				return
			}
			log.Printf("AuthMiddleware ERROR for session %s: %+v", c.Value, err)
			httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
			return
		}

		// Setting up context
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextRequesterAccountIDKey, s.AccountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
