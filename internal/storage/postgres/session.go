package postgres

import (
	"fmt"
	"time"

	"github.com/dmartzol/hackerspace/internal/models"
	_ "github.com/lib/pq"
)

// SessionFromIdentifier fetches a session from a given its identifier
func (db *DB) SessionFromIdentifier(identifier string) (*models.Session, error) {
	var s models.Session
	sqlStatement := `select * from sessions where session_id = $1`
	err := db.Get(&s, sqlStatement, identifier)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSession creates a new session
func (db *DB) CreateSession(accountID int64) (*models.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var s models.Session
	sqlStatement := `insert into sessions (account_id) values ($1) returning *`
	err = tx.Get(&s, sqlStatement, accountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
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
func (db *DB) UpdateSession(identifier string) (*models.Session, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var s models.Session
	sqlStatement := `update sessions set last_activity_time=default where session_id = $1 returning *`
	err = tx.Get(&s, sqlStatement, identifier)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
}
