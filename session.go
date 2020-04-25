package main

import "time"

// Session represents an account session
type Session struct {
	Row
	AccountID      int64
	ExpirationDate time.Time
	Token          string
}
