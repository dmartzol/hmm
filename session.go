package main

import "time"

const (
	// sessionLength represents the duration(in minutes) a session will be valid for
	sessionLength = 120
)

// Session represents an account session
type Session struct {
	Row
	AccountID         int64     `db:"account_id"`
	LastActivity      time.Time `db:"last_activity"`
	SessionIdentifier string    `db:"session_id"`
}
