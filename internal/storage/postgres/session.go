package postgres

import (
	"log"
	"time"

	"github.com/dmartzol/hmm/internal/domain"
	_ "github.com/lib/pq"
)

var (
	ErrExpiredResource error
)

// SessionFromToken fetches a session by its token
func (db *DB) SessionFromToken(token string) (*domain.Session, error) {
	var s domain.Session
	sqlStatement := `select * from sessions where token = $1`
	err := db.Get(&s, sqlStatement, token)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSession creates a new session
func (db *DB) CreateSession(accountID int64) (*domain.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var s domain.Session
	sqlStatement := `insert into sessions (account_id) values ($1) returning *`
	err = tx.Get(&s, sqlStatement, accountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
}

// ExpireSessionFromToken expires the session with the given token
func (db *DB) ExpireSessionFromToken(token string) (*domain.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var s domain.Session
	sqlStatement := `update sessions set expiration_time = current_timestamp where token = $1 returning *`
	err = tx.Get(&s, sqlStatement, token)
	if err != nil {
		tx.Rollback()
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
func (db *DB) UpdateSession(token string) (*domain.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var session domain.Session
	sqlStatement := `select * from sessions where token = $1`
	tx.Get(&session, sqlStatement, token)
	if err != nil {
		log.Printf("UpdateSession db - ERROR fetching session from token %s: %+v", token, err)
		tx.Rollback()
		return nil, err
	}
	if session.ExpirationTime.Before(time.Now()) {
		return nil, ErrExpiredResource
	}
	var updatedSession domain.Session
	sqlStatement = `update sessions set last_activity_time=default where token = $1 returning *`
	err = tx.Get(&updatedSession, sqlStatement, token)
	if err != nil {
		log.Printf("UpdateSession db - ERROR updating session from token %s: %+v", token, err)
		tx.Rollback()
		return nil, err
	}
	return &updatedSession, tx.Commit()
}
