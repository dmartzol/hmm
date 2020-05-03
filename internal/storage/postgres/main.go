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
	ID         int64
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
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
