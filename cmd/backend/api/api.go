package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/dmartzol/go-sdk/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

const (
	apiVersionNumber = "0.0.1"
	idQueryParameter = "id"
	appName          = "backend_api"
)

type API struct {
	http.Handler
	*Resources
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewAPI(db *sqlx.DB, logger logger.Logger) *API {
	resources := newResources(db, logger)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	//r.Use(h.AuthMiddleware)

	r.Get("/version", resources.Version)

	// accounts
	r.Post("/accounts", resources.CreateAccount)
	//r.HandleFunc("/accounts/{id}", h.GetAccount).Methods("GET")
	//r.HandleFunc("/accounts", h.GetAccounts).Methods("GET")
	//r.HandleFunc("/accounts/{id}/confirm-email", h.ConfirmEmail).Methods("POST")
	//r.HandleFunc("/accounts/password", h.ResetPassword).Methods("POST")
	//r.HandleFunc("/accounts/{id}/roles", h.AddAccountRole).Methods("POST")
	//r.HandleFunc("/accounts/{id}/roles", h.GetAccountRoles).Methods("GET")

	// RESTy routes for "accounts" resource
	//router.Route("/accounts", func(r chi.Router) {
	////r.With(paginate).Get("/", ListArticles)
	//r.Post("/", http.HandlerFunc(notImplementedHandler))      // POST /accounts
	//r.Get("/search", http.HandlerFunc(notImplementedHandler)) // GET /accounts/search

	//router.Route("/{accountID}", func(r chi.Router) {
	//r.Use(AccountCtx)                                      // Load the *Account on the request context
	//r.Get("/", http.HandlerFunc(notImplementedHandler))    // GET /accounts/123
	//r.Put("/", http.HandlerFunc(notImplementedHandler))    // PUT /accounts/123
	//r.Delete("/", http.HandlerFunc(notImplementedHandler)) // DELETE /accounts/123
	//})

	//// GET /accounts/whats-up
	//router.With(AccountCtx).Get("/{accountSlug:[a-z-]+}", nil)
	//})

	//r.HandleFunc("/users", APIUsers).Methods("GET", "POST")

	//// return 405 for PUT, PATCH and DELETE
	//r.HandleFunc("/users", status(405, "GET", "POST")).Methods("PUT", "PATCH", "DELETE")

	//r = r.PathPrefix("/v1").Subrouter()

	//r.HandleFunc("/version", h.Version).Methods("GET")

	//// sessions
	//// see: https://stackoverflow.com/questions/7140074/restfully-design-login-or-register-resources
	//r.HandleFunc("/sessions", h.CreateSession).Methods("POST")
	//r.HandleFunc("/sessions", h.GetSession).Methods("GET")
	//r.HandleFunc("/sessions", h.ExpireSession).Methods("DELETE")

	//// roles
	//r.HandleFunc("/roles", h.GetRoles).Methods("GET")
	//r.HandleFunc("/roles", h.CreateRole).Methods("POST")
	//r.HandleFunc("/roles/{id}", h.EditRole).Methods("PUT")

	//cors := cors.New(cors.Options{
	//AllowedOrigins:   []string{"http://localhost:3000"},
	//AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	//AllowCredentials: true,
	//// Enable Debugging for testing, consider disabling in production
	//// Debug: true,
	//})

	return &API{
		Resources: resources,
		Handler:   r,
	}
}
