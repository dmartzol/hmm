package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/pkg/timeutils"
)

// Account is the restricted response body of hmm.Account
// see: https://stackoverflow.com/questions/46427723/golang-elegant-way-to-omit-a-json-property-from-being-serialized
type Account struct {
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
	Roles                      []Role
}

func AccountView(a *hmm.Account, options map[string]bool) Account {
	view := Account{
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
			view.Roles = append(view.Roles, RoleView(r, nil))
		}
	}
	return view
}

func AccountsView(accs hmm.Accounts, options map[string]bool) []Account {
	var l []Account
	for _, a := range accs {
		l = append(l, AccountView(a, options))
	}
	return l
}

type RegisterRequest struct {
	FirstName   string
	LastName    string
	DOB         string
	DOBTime     time.Time
	Gender      *string
	PhoneNumber *string
	Email       string
	Password    string
}

func (r *RegisterRequest) ValidateAndNormalize() error {
	err := r.validate()
	if err != nil {
		return fmt.Errorf("error validating: %w", err)
	}
	return nil
}

func (r RegisterRequest) validate() error {
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
	return nil
}

func (r *RegisterRequest) normalize() error {
	r.FirstName = NormalizeName(r.FirstName)
	r.LastName = NormalizeName(r.LastName)
	var err error
	r.DOBTime, err = time.Parse(timeutils.LayoutISODay, r.DOB)
	if err != nil {
		return fmt.Errorf("error parsing DOB %q: %w", r.DOB, err)
	}
	return nil
}

func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	return name
}
