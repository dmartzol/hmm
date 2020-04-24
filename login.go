package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type loginCredentials struct {
	Email    string
	Password string
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var credentials loginCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// fetching account with credentials(errors reurned should be purposedly broad)
	registered, err := db.EmailExists(credentials.Email)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if !registered {
		log.Printf("unable to find email '%s' in db", credentials.Email)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	a, err := db.AccountWithCredentials(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	credentials.Password = ""

	// create session and cookie
	s, err := db.CreateSession(a.ID, time.Now().Add(time.Minute*5), uuid.New().String())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "session",
			Value: s.Token,
		}
		http.SetCookie(w, cookie)
	}

	json.NewEncoder(w).Encode(s)
}

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	s, err := db.SessionFromToken(c.Value)
	if err != nil {
		return false
	}
	if s.ExpirationDate.Before(time.Now()) {
		return false
	}
	return true
}

func logout(w http.ResponseWriter, r *http.Request) {

}
