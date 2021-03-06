package models

import (
	"errors"
	"strings"
	"time"

	"github.com/dmartzol/hmm/pkg/timeutils"
)

type Accounts []*Account

// Account represents a user account
type Account struct {
	Row
	FirstName                 string `db:"first_name"`
	LastName                  string `db:"last_name"`
	DOB                       time.Time
	Gender                    *string
	Active                    bool
	FailedLoginsCount         int64   `db:"failed_logins_count"`
	DoorCode                  *string `db:"door_code"`
	PassHash                  string
	Email                     string
	ConfirmedEmail            bool       `db:"confirmed_email"`
	PhoneNumber               *string    `db:"phone_number"`
	ConfirmedPhone            bool       `db:"confirmed_phone"`
	ZipCode                   string     `db:"zip_code"`
	ReviewTime                *time.Time `db:"review_time"` // timestamp of when the account was reviewed
	ExternalPaymentCustomerID *int64     `db:"external_payment_customer_id"`

	// fields to populate
	Roles Roles
}

// AccountView is the restricted response body of Account
// see: https://stackoverflow.com/questions/46427723/golang-elegant-way-to-omit-a-json-property-from-being-serialized
type AccountView struct {
	ID                         int64 `json:"ID"`
	FirstName, LastName, Email string
	DOB                        string `json:"DateOfBird"`
	PhoneNumber                string `json:",omitempty"`
	DoorCode                   string `json:",omitempty"`
	Gender                     string `json:",omitempty"`
	Active                     bool
	ConfirmedEmail             bool
	ConfirmedPhone             bool
	FailedLoginsCount          int64
	Roles                      []RoleView
}

// View returns the Account struct restricted to those fields allowed in options
// see: https://stackoverflow.com/questions/46427723/golang-elegant-way-to-omit-a-json-property-from-being-serialized
func (a Account) View(options map[string]bool) AccountView {
	view := AccountView{
		ID:                a.ID,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		DOB:               a.DOB.Format(timeutils.LayoutISODay),
		Active:            a.Active,
		FailedLoginsCount: a.FailedLoginsCount,
		Email:             a.Email,
		ConfirmedEmail:    a.ConfirmedEmail,
		ConfirmedPhone:    a.ConfirmedPhone,
	}
	if a.DoorCode != nil && options["door_code"] {
		view.DoorCode = *a.DoorCode
	}
	if a.PhoneNumber != nil && options["phone_number"] {
		view.PhoneNumber = *a.PhoneNumber
	}
	if a.Gender != nil {
		view.Gender = *a.Gender
	}
	if a.Roles != nil {
		for _, r := range a.Roles {
			view.Roles = append(view.Roles, r.View(nil))
		}
	}
	return view
}

func (accs Accounts) Views(options map[string]bool) []AccountView {
	var l []AccountView
	for _, a := range accs {
		l = append(l, a.View(options))
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
	if r.Gender != nil && *r.Gender != "" && *r.Gender != "M" && *r.Gender != "F" {
		return errors.New("Gender value not implemented")
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
