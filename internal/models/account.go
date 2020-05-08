package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dmartzol/hackerspace/pkg/timeutils"
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
	FailedLoginsCount int64   `db:"failed_logins_count"`
	DoorCode          *string `db:"door_code"`
	PassHash          string
	Email             string
	PhoneNumber       *string `db:"phone_number"`

	RoleID                    *int64 `db:"role_id"`
	ExternalPaymentCustomerID *int64 `db:"external_payment_customer_id"`
}

// AccountAPI should always be the object used to respond to any request
type AccountAPI struct {
	Row
	FirstName, LastName, Email    string
	DOB                           time.Time `json:"DateOfBird"`
	Gender, DoorCode, PhoneNumber string    `json:"omitempty"`
	Active                        bool
	FailedLoginsCount             int64
}

func (a Account) API() AccountAPI {
	r := AccountAPI{
		Row:               a.Row,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		DOB:               a.DOB,
		Active:            a.Active,
		FailedLoginsCount: a.FailedLoginsCount,
		Email:             a.Email,
	}
	if a.DoorCode != nil {
		r.DoorCode = *a.DoorCode
	}
	if a.PhoneNumber != nil {
		r.PhoneNumber = *a.PhoneNumber
	}
	if a.Gender != nil {
		r.Gender = *a.Gender
	}
	return r
}

func (accs Accounts) API() []AccountAPI {
	l := []AccountAPI{}
	for _, a := range accs {
		l = append(l, a.API())
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
		return errors.New("must provide first name")
	}
	if r.LastName == "" {
		return errors.New("must provide last name")
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
	ConfirmationCode string
}
