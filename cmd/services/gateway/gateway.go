package main

import (
	"log"
	"net/http"

	"github.com/dmartzol/hmm/internal/handler"
	"github.com/dmartzol/hmm/internal/storage/postgres"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/cli"
)

func newGatewayServiceRun(c *cli.Context) error {
	port := c.String(flagPort)
	host := c.String(flagHost)

	dbConfig := postgres.Config{
		Host:     c.String(flagDBHost),
		Port:     c.Int(flagDBPort),
		Name:     c.String(flagDBName),
		User:     c.String(flagDBUser),
		Password: c.String(flagDBPass),
	}
	db, err := postgres.New(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	h, err := handler.New(db)
	if err != nil {
		log.Fatalf("error starting h: %+v", err)
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/version", h.Version).Methods("GET")

	// sessions
	// see: https://stackoverflow.com/questions/7140074/restfully-design-login-or-register-resources
	r.HandleFunc("/sessions", h.CreateSession).Methods("POST")
	r.HandleFunc("/sessions", h.GetSession).Methods("GET")
	r.HandleFunc("/sessions", h.ExpireSession).Methods("DELETE")

	// accounts
	r.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", h.GetAccount).Methods("GET")
	r.HandleFunc("/accounts", h.GetAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}/confirm-email", h.ConfirmEmail).Methods("POST")
	r.HandleFunc("/accounts/password", h.ResetPassword).Methods("POST")
	r.HandleFunc("/accounts/{id}/roles", h.AddAccountRole).Methods("POST")
	r.HandleFunc("/accounts/{id}/roles", h.GetAccountRoles).Methods("GET")

	// roles
	r.HandleFunc("/roles", h.GetRoles).Methods("GET")
	r.HandleFunc("/roles", h.CreateRole).Methods("POST")
	r.HandleFunc("/roles/{id}", h.EditRole).Methods("PUT")

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		h.AuthMiddleware,
	)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		// Debug: true,
	})

	log.Print("listening and serving")
	return http.ListenAndServe(host+":"+port, cors.Handler(r))
}
