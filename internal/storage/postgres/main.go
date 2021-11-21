package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB represents the database
type DB struct {
	*sqlx.DB
}

type Config struct {
	Name, User, Password, Host string
	Port                       int
}

func NewDB(config Config) (*DB, error) {
	var dataSourceName string
	dataSourceName = "host=%s port=%d user=%s dbname=%s sslmode=disable"
	dataSourceName = fmt.Sprintf(dataSourceName, config.Host, config.Port, config.User, config.Name)
	if config.Password != "" {
		dataSourceName += fmt.Sprintf(" password=%s", config.Password)
	}
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
	return &DB{database}, nil
}
