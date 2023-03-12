package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dmartzol/hmm/internal/hmm"
	_ "github.com/lib/pq"
)

var (
	ErrExpiredResource    error
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// SessionFromToken fetches a session by its token
func (db *DB) SessionFromToken(token string) (*hmm.Session, error) {
	var s hmm.Session
	sqlStatement := `select * from sessions where token = $1`
	err := db.Get(&s, sqlStatement, token)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSession creates a new session
func (db *DB) CreateSession(email, password string) (*hmm.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	var a hmm.Account
	sqlSelect := `select * from accounts a where a.email = $1 and a.passhash = crypt($2, a.passhash)`
	err = tx.Get(&a, sqlSelect, email, password)
	if err == sql.ErrNoRows {
		_ = tx.Rollback()
		return nil, ErrInvalidCredentials
	} else if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("error fetching account for email %q: %w", email, err)
	}

	var s hmm.Session
	sqlInsert := `insert into sessions (account_id) values ($1) returning *`
	err = tx.Get(&s, sqlInsert, a.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("error creating session for account %q: %w", a.ID, err)
	}
	return &s, tx.Commit()
}

// ExpireSession expires the session with the given token
func (db *DB) ExpireSession(token string) (*hmm.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var s hmm.Session
	sqlStatement := `update sessions set expiration_time = current_timestamp where token = $1 returning *`
	err = tx.Get(&s, sqlStatement, token)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
}

// CleanSessionsOlderThan deletes all sessions older than age(in seconds) and returns the number of rows affected
func (db *DB) CleanSessionsOlderThan(age time.Duration) (int64, error) {
	t := time.Now().Add(-age * time.Second)
	sqlStatement := `delete from sessions where last_activity_time < $1`
	res, err := db.Exec(sqlStatement, t)
	if err != nil {
		return -1, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return count, nil
}

// UpdateSession updates a session in the db with the current timestamp
func (db *DB) UpdateSession(token string) (*hmm.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var session hmm.Session
	sqlStatement := `select * from sessions where token = $1`
	err = tx.Get(&session, sqlStatement, token)
	if err != nil {
		log.Printf("UpdateSession db - ERROR fetching session from token %s: %+v", token, err)
		_ = tx.Rollback()
		return nil, err
	}
	if session.ExpirationTime.Before(time.Now()) {
		return nil, ErrExpiredResource
	}
	var updatedSession hmm.Session
	sqlStatement = `update sessions set last_activity_time=default where token = $1 returning *`
	err = tx.Get(&updatedSession, sqlStatement, token)
	if err != nil {
		log.Printf("UpdateSession db - ERROR updating session from token %s: %+v", token, err)
		_ = tx.Rollback()
		return nil, err
	}
	return &updatedSession, tx.Commit()
}
