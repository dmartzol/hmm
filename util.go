package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
)

// GetEnvString get an environment variable and if it doesn't exists, it returns fallback
func GetEnvString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvInt get an environment variable and if it doesn't exists, it returns fallback
func GetEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	p, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("error parsing integer from env variable '%s': %+v", key, err)
		return fallback
	}
	return p
}

func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	if len(*s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func NewNullInt(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}
