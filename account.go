package main

import (
	"net/http"
	"time"
)

// Account represents a user account
type Account struct {
	Row
	FirstName           string
	LastName            string
	Suffix              string
	DOB                 time.Time
	Gender              *string
	Active              bool
	Email               string
	FailedLoginAttempts int64
	DoorCode            string

	RoleID                    int64
	PhoneNumberID             int64
	ExternalPaymentCustomerID int64
}

func register(w http.ResponseWriter, r *http.Request) {

}
