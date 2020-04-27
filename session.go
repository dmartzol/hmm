package main

import "time"

const (
	// SessionLimit represents the duration(in minutes) a session will be valid for
	SessionLimit = 5
)

// Session represents an account session
type Session struct {
	Row
	AccountID      int64     `db:"account_id"`
	ExpirationDate time.Time `db:"expiration_date"`
	Token          string
}
