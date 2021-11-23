package hmm

import (
	"strings"
	"time"
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

type AccountService interface {
	Create(a *Account, password, confirmationCode string) (*Account, *Confirmation, error)
	Account(id int64) (*Account, error)
	Accounts() (Accounts, error)
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