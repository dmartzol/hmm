package postgres

import (
	"os"

	"github.com/dmartzol/hackerspace/pkg/environment"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB represents the database
type DB struct {
	*sqlx.DB
}

const (
	dbport = "DBPORT"
	dbuser = "PGUSER"
	dbpass = "PGPASSWORD"
	dbhost = "PGHOST"
	dbname = "PGDATABASE"
)

type DatabaseConfig struct {
	Name, User, Password, Host string
	Port                       int
}

func DBConfig() DatabaseConfig {
	config := DatabaseConfig{}
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
	config.Port = environment.GetEnvInt(dbport, 5432)
	config.Password = environment.GetEnvString(dbpass, "")
	config.Host = host
	config.User = user
	config.Name = name
	return config
}
