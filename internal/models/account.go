package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dmartzol/hmmm/pkg/timeutils"
)

type Accounts []*Account

// Account represents a user account
type Account struct {
	Row
	FirstName         string `db:"first_name"`
	LastName          string `db:"last_name"`
	DOB               time.Time
	Gender            *string
	Active            bool
	Reviewed          bool
	FailedLoginsCount int64   `db:"failed_logins_count"`
	DoorCode          *string `db:"door_code"`
	PassHash          string
	Email             string
	ConfirmedEmail    bool    `db:"confirmed_email"`
	PhoneNumber       *string `db:"phone_number"`
	ConfirmedPhone    bool    `db:"confirmed_phone"`
	ZipCode           string  `db:"zip_code"`

	RoleID                    *int64 `db:"role_id"`
	ExternalPaymentCustomerID *int64 `db:"external_payment_customer_id"`
}

// AccountView is the restricted response body of Account
// see: https://stackoverflow.com/questions/46427723/golang-elegant-way-to-omit-a-json-property-from-being-serialized
type AccountView struct {
	Row
	FirstName, LastName, Email string
	DOB                        time.Time `json:"DateOfBird"`
	PhoneNumber                string    `json:",omitempty"`
	DoorCode                   string    `json:",omitempty"`
	Gender                     string    `json:",omitempty"`
	Active                     bool
	FailedLoginsCount          int64
}

// Restrict returns the Account struct restricted to those fields allowed in options
// see: https://stackoverflow.com/questions/46427723/golang-elegant-way-to-omit-a-json-property-from-being-serialized
func (a Account) Restrict(options map[string]bool) AccountView {
	r := AccountView{
		Row:               a.Row,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		DOB:               a.DOB,
		Active:            a.Active,
		FailedLoginsCount: a.FailedLoginsCount,
		Email:             a.Email,
	}
	if a.DoorCode != nil && options["door_code"] {
		r.DoorCode = *a.DoorCode
	}
	if a.PhoneNumber != nil && options["phone_number"] {
		r.PhoneNumber = *a.PhoneNumber
	}
	if a.Gender != nil {
		r.Gender = *a.Gender
	}
	return r
}

func (accs Accounts) Restrict(options map[string]bool) []AccountView {
	l := []AccountView{}
	for _, a := range accs {
		l = append(l, a.Restrict(options))
	}
	return l
}

type RegisterRequest struct {
	FirstName   string
	LastName    string
	DOB         string
	Gender      *string
	PhoneNumber *string
	Email       string
	Password    string
}

func (r RegisterRequest) Validate() error {
	if r.FirstName == "" {
		return errors.New("first name is required")
	}
	if r.LastName == "" {
		return errors.New("last name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password too short")
	}
	_, err := time.Parse(timeutils.LayoutISO, r.DOB)
	if err != nil {
		return fmt.Errorf("time.Parse %v: %w", r.DOB, err)
	}
	return nil
}

func validEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	if !strings.Contains(email, ".") {
		return false
	}
	return true
}

type ResetPasswordRequest struct {
	Email string
}

type ConfirmEmailRequest struct {
	ConfirmationKey string
}
