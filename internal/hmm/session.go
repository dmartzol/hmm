package hmm

import (
	"time"
)

type SessionService interface {
	Create(email, password string) (*Session, error)
	SessionFromToken(token string) (*Session, error)
	ExpireSession(token string) (*Session, error)
	UpdateSession(token string) (*Session, error)
}

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
	ExpirationTime   time.Time
}

func (s Session) View(options map[string]bool) SessionView {
	view := SessionView{
		AccountID:        s.AccountID,
		LastActivityTime: s.LastActivityTime,
		ExpirationTime:   s.ExpirationTime,
	}
	return view
}

type LoginCredentials struct {
	Email    string
	Password string
}
