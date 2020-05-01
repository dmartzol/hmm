package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

const (
	apiVersion            = "0.0.1"
	hackerSpaceCookieName = "HackerSpace-Cookie"
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
	fmt.Fprintf(w, "Home")
}

func main() {
	log.SetFlags(LstdFlags)

	r := mux.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(authMiddleware)

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/version", version).Methods("GET")

	// sessions
	// see: https://stackoverflow.com/questions/7140074/restfully-design-login-or-register-resources
	r.HandleFunc("/sessions", login).Methods("POST")
	r.HandleFunc("/sessions/{id}", logout).Methods("DELETE")

	// accounts
	r.HandleFunc("/accounts", register).Methods("POST")
	r.HandleFunc("/accounts/{id}", getAccount).Methods("GET")

	log.Print("listening and serving")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ss := map[string]string{
			"/version":  "GET",
			"/sessions": "POST",
			"/accounts": "POST",
		}
		method, in := ss[r.RequestURI]
		if in && method == r.Method {
			next.ServeHTTP(w, r)
			return
		}
		c, err := r.Cookie(hackerSpaceCookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			log.Printf("cookie: %+v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		_, err = db.UpdateSession(c.Value)
		if err != nil {
			log.Printf("UpdateSession: %+v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
