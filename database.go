package main

import (
	"fmt"
	"log"
	"os"
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

// Row represents a database row
type Row struct {
	ID        int64
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func init() {
	dbConfig := dbConfig()

	dataSourceName := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dataSourceName = fmt.Sprintf(dataSourceName, dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
	database, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}
	err = database.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	db = &DB{database}
}

type databaseConfig struct {
	Name, User, Password, Host string
	Port                       int
}

func dbConfig() databaseConfig {
	config := databaseConfig{}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("PGDATABASE environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("PGUSER environment variable required but not set")
	}
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("PGHOST environment variable required but not set")
	}
	config.Port = GetEnvInt(dbport, 5432)
	config.Password = GetEnvString(dbpass, "")
	config.Host = host
	config.User = user
	config.Name = name
	return config
}

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
	sqlStatement := `delete from sessions where last_activity < $1`
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
	sqlStatement := `update sessions set last_activity=default where session_id = $1 returning *`
	err = tx.Get(&s, sqlStatement, identifier)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &s, tx.Commit()
}

// EmailExists returns true if the provided email exists in the db
func (db *DB) EmailExists(email string) (bool, error) {
	var exists bool
	sqlStatement := `select exists(select 1 from accounts a where a.email = $1)`
	err := db.Get(&exists, sqlStatement, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// AccountWithCredentials returns an account if the email and password provided match an (email,password) pair in the db
func (db *DB) AccountWithCredentials(email, allegedPassword string) (*Account, error) {
	var a Account
	sqlStatement := `select * from accounts a where a.email = $1 and a.passhash = crypt($2, a.passhash)`
	err := db.Get(&a, sqlStatement, email, allegedPassword)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CreateAccount creates a new account in the db
func (db *DB) CreateAccount(first, last, email, password string, dob time.Time, gender, phone *string) (*Account, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var a Account
	sqlStatement := `insert into accounts (first_name, last_name, dob, gender, phone_number, email, passhash) values ($1, $2, $3, $4, $5, $6, crypt($7, gen_salt('bf', 8))) returning *`
	err = tx.Get(&a, sqlStatement, first, last, dob, gender, phone, email, password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &a, tx.Commit()
}
