package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	apiVersion = "0.0.1"
)

type loginRequest struct {
	Email    string
	Password string
}

func version(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "version %s", apiVersion)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome")
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", index)
	r.Get("/version", version)

	// sessions
	// see: https://stackoverflow.com/questions/7140074/restfully-design-login-or-register-resources
	r.Post("/session", login)
	r.Delete("/session", logout)

	// accounts
	r.Post("/account", register)

	log.Fatal(http.ListenAndServe("localhost:8080", r))

}
