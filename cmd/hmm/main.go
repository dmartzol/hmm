package main

import (
	"log"
	"net/http"

	"github.com/dmartzol/hmm/internal/controllers"
	"github.com/dmartzol/hmm/internal/storage/postgres"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	Ldate         = 1 << iota                  // the date in the local time zone: 2009/01/23
	Ltime                                      // the time in the local time zone: 01:23:23
	Lmicroseconds                              // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                                  // full file name and line number: /a/b/c/d.go:23
	Lshortfile                                 // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                       // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                                 // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime | Lshortfile // initial values for the standard logger
)

func main() {
	log.SetFlags(LstdFlags)
	db, err := postgres.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	api, err := controllers.NewAPI(db)
	if err != nil {
		log.Fatalf("error starting api: %+v", err)
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/version", api.Version).Methods("GET")

	// sessions
	// see: https://stackoverflow.com/questions/7140074/restfully-design-login-or-register-resources
	r.HandleFunc("/sessions", api.CreateSession).Methods("POST")
	r.HandleFunc("/sessions", api.GetSession).Methods("GET")
	r.HandleFunc("/sessions", api.ExpireSession).Methods("DELETE")

	// accounts
	r.HandleFunc("/accounts", api.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", api.GetAccount).Methods("GET")
	r.HandleFunc("/accounts", api.GetAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}/confirm-email", api.ConfirmEmail).Methods("POST")
	r.HandleFunc("/accounts/password", api.ResetPassword).Methods("POST")
	r.HandleFunc("/accounts/{id}/roles", api.AddAccountRole).Methods("POST")
	r.HandleFunc("/accounts/{id}/roles", api.GetAccountRoles).Methods("GET")

	// roles
	r.HandleFunc("/roles", api.GetRoles).Methods("GET")
	r.HandleFunc("/roles", api.CreateRole).Methods("POST")
	r.HandleFunc("/roles/{id}", api.EditRole).Methods("PUT")

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		api.AuthMiddleware,
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		// Debug: true,
	})

	log.Print("listening and serving")
	log.Fatal(http.ListenAndServe("localhost:3001", c.Handler(r)))
}
