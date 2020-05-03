package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB represents the database
type DB struct {
	*sqlx.DB
}

var db *DB

const (
	dbport = "DBPORT"
	dbuser = "PGUSER"
	dbpass = "PGPASSWORD"
	dbhost = "PGHOST"
	dbname = "PGDATABASE"
)

// SessionFromIdentifier fetches a session from a given its identifier
func (db *DB) SessionFromIdentifier(identifier string) (*Session, error) {
	var s Session
	sqlStatement := `select * from sessions where session_id = $1`
	err := db.Get(&s, sqlStatement, identifier)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSession creates a new session
func (db *DB) CreateSession(accountID int64) (*Session, error) {
	var s Session
	sqlStatement := `insert into sessions (account_id) values ($1) returning *`
	err := db.Get(&s, sqlStatement, accountID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// DeleteSession deletes the session with the given identifier
func (db *DB) DeleteSession(identifier string) error {
	sqlStatement := `delete from sessions where session_id = $1`
	res, err := db.Exec(sqlStatement, identifier)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("%d rows affected by DeleteSession", count)
	}
	fmt.Println(count)
	return nil
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

// UpdateSession sets the current timestamp
func (db *DB) UpdateSession(identifier string) (*Session, error) {
	tx, err := db.Beginx()
	var s Session
	sqlStatement := `update sessions set last_activity_time=default where session_id = $1 returning *`
	err = tx.Get(&s, sqlStatement, identifier)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
}
