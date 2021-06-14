package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/dmartzol/hmm/internal/models"
	_ "github.com/lib/pq"
)

// SessionFromToken fetches a session by its token
func (db *DB) SessionFromToken(token string) (*models.Session, error) {
	var s models.Session
	sqlStatement := `select * from sessions where token = $1`
	err := db.Get(&s, sqlStatement, token)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Login creates a new session from credentials request
func (db *DB) Login(r models.LoginCredentials) (*models.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	a, err := db.AccountFromCredentialsUsingTx(tx, r.Email, r.Password)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, ErrResourceDoesNotExist
		}
		return nil, err
	}
	var s models.Session
	sqlStatement := `insert into sessions (account_id) values ($1) returning *`
	err = tx.Get(&s, sqlStatement, a.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
}

// CreateSession creates a new session from account id
func (db *DB) CreateSession(accountID int64) (*models.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	a, err := db.AccountUsingTx(tx, accountID)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, ErrResourceDoesNotExist
		}
		return nil, err
	}
	var s models.Session
	sqlStatement := `insert into sessions (account_id) values ($1) returning *`
	err = tx.Get(&s, sqlStatement, a.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
}

// ExpireSessionFromToken expires the session with the given token
func (db *DB) ExpireSessionFromToken(token string) (*models.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var s models.Session
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
func (db *DB) UpdateSession(token string) (*models.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var session models.Session
	sqlStatement := `select * from sessions where token = $1`
	err = tx.Get(&session, sqlStatement, token)
	if err != nil {
		log.Printf("UpdateSession db - ERROR fetching session from token %s: %+v", token, err)
		tx.Rollback()
		return nil, err
	}
	if session.ExpirationTime.Before(time.Now()) {
		return nil, ErrExpiredResource
	}
	var updatedSession models.Session
	sqlStatement = `update sessions set last_activity_time=default where token = $1 returning *`
	err = tx.Get(&updatedSession, sqlStatement, token)
	if err != nil {
		log.Printf("UpdateSession db - ERROR updating session from token %s: %+v", token, err)
		tx.Rollback()
		return nil, err
	}
	return &updatedSession, tx.Commit()
}
