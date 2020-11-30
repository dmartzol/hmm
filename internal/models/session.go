package models

import (
	"time"
)

// Session represents an account session
type Session struct {
	Row
	AccountID        int64     `db:"account_id"`
	Token            string    `db:"token"`
	LastActivityTime time.Time `db:"last_activity_time"`
	ExpirationTime   time.Time `db:"expiration_time"`
}

type SessionView struct {
	AccountID        int64
	LastActivityTime time.Time
}

func (s Session) View(options map[string]bool) SessionView {
	view := SessionView{
		AccountID:        s.AccountID,
		LastActivityTime: s.LastActivityTime,
	}
	return view
}

type LoginCredentials struct {
	Email    string
	Password string
}
