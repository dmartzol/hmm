package main

import (
	"log"
	"net/http"

	"github.com/dmartzol/hackerspace/internal/storage/postgres"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
)

const (
	apiVersionNumber      = "0.0.1"
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

type storage interface {
	sessionStorage
	accountStorage
}

// API represents something
type API struct {
	storage
}

func newAPI() (*API, error) {
	db, err := postgres.NewDB()
	if err != nil {
		return nil, err
	}
	return &API{db}, nil
}

func main() {
	log.SetFlags(LstdFlags)
	api, err := newAPI()
	if err != nil {
		log.Fatalf("error starting api: %+v", err)
	}

	r := mux.NewRouter()
	r = r.PathPrefix("/v1").Subrouter()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		api.authMiddleware,
	)

	r.HandleFunc("/version", api.version).Methods("GET")

	// sessions
	// see: https://stackoverflow.com/questions/7140074/restfully-design-login-or-register-resources
	r.HandleFunc("/sessions", api.createSession).Methods("POST")
	r.HandleFunc("/sessions", api.deleteSession).Methods("DELETE")

	// accounts
	r.HandleFunc("/accounts", api.createAccount).Methods("POST")
	r.HandleFunc("/accounts/{id}", api.getAccount).Methods("GET")
	r.HandleFunc("/accounts", api.getAccounts).Methods("GET")
	r.HandleFunc("/accounts/{id}/confirm-email", api.confirmEmail).Methods("POST")
	r.HandleFunc("/accounts/password", api.resetPassword).Methods("POST")

	log.Print("listening and serving")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func (api API) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ss := map[string]string{
			"/v1/version":  "GET",
			"/v1/sessions": "POST",
			"/v1/accounts": "POST",
		}
		method, in := ss[r.RequestURI]
		if in && method == r.Method {
			next.ServeHTTP(w, r)
			return
		}
		c, err := r.Cookie(hackerSpaceCookieName)
		if err != nil {
			if err != http.ErrNoCookie {
				log.Printf("cookie: %+v", err)
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		_, err = api.UpdateSession(c.Value)
		if err != nil {
			log.Printf("UpdateSession: %+v", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
