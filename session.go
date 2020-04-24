package main

import "time"

type Session struct {
	Row
	AccountID      int64
	ExpirationDate time.Time
	Token          string
}
