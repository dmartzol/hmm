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
	log.SetFlags(LstdFlags)

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

	log.Print("ListenAndServe")
	log.Fatal(http.ListenAndServe("localhost:8080", r))

}
