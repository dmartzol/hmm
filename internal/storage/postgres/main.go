package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/dmartzol/hackerspace/pkg/environment"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	dbport = "DBPORT"
	dbuser = "PGUSER"
	dbpass = "PGPASSWORD"
	dbhost = "PGHOST"
	dbname = "PGDATABASE"
)

// DB represents the database
type DB struct {
	*sqlx.DB
}

type DatabaseConfig struct {
	Name, User, Password, Host string
	Port                       int
}

func (db *DB) PrepareDatabase() (*sqlx.DB, error) {
	dbConfig := DatabaseConfig{}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		return nil, fmt.Errorf("PGDATABASE environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		return nil, fmt.Errorf("PGUSER environment variable required but not set")
	}
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		return nil, fmt.Errorf("PGHOST environment variable required but not set")
	}
	dbConfig.Port = environment.GetEnvInt(dbport, 5432)
	dbConfig.Password = environment.GetEnvString(dbpass, "")
	dbConfig.Host = host
	dbConfig.User = user
	dbConfig.Name = name

	dataSourceName := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	dataSourceName = fmt.Sprintf(dataSourceName, dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
	database, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Printf("error connecting to db: %+v", err)
		return nil, err
	}
	err = database.Ping()
	if err != nil {
		log.Printf("error pinging db: %+v", err)
		return nil, err
	}
	return database, nil
}
