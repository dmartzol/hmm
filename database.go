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

// SessionFromToken fetches a session from a given token
func (db *DB) SessionFromToken(token string) (*Session, error) {
	var s Session
	sqlQuery := `select * from sessions where token = $1`
	err := db.Get(&s, sqlQuery, token)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSession creates a new session
func (db *DB) CreateSession(accountID int64, expiry time.Time, token string) (*Session, error) {
	var s Session
	sqlQuery := `insert into sessions (account_id, expiration_date, token) values ($1, $2, $3) returning *`
	err := db.Get(&s, sqlQuery, accountID, expiry, token)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// EmailExists returns true if the provided email exists in the db
func (db *DB) EmailExists(email string) (bool, error) {
	var exists bool
	sqlQuery := `select exists(select 1 from accounts a where a.email = $1)`
	err := db.Get(&exists, sqlQuery, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// AccountWithCredentials returns an account if the email and password provided match an (email,password) pair in the db
func (db *DB) AccountWithCredentials(email, allegedPassword string) (*Account, error) {
	var a Account
	sqlQuery := `select * from accounts a where a.email = $1 and a.passhash = crypt($2, a.passhash)`
	err := db.Get(&a, sqlQuery, email, allegedPassword)
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
	sqlQuery := `insert into accounts (first_name, last_name, dob, gender, phone_number, email, passhash) values ($1, $2, $3, $4, $5, $6, crypt($7, gen_salt('bf', 8))) returning *`
	err = tx.Get(&a, sqlQuery, first, last, dob, gender, phone, email, password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &a, tx.Commit()
}
