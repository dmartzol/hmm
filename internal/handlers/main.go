package handlers

import (
	"github.com/dmartzol/hackerspace/internal/storage/postgres"
)

const (
	apiVersionNumber      = "0.0.1"
	hackerSpaceCookieName = "HackerSpace-Cookie"
)

type storage interface {
	sessionStorage
	accountStorage
}

// API represents something
type API struct {
	storage
}

func NewAPI() (*API, error) {
	db, err := postgres.NewDB()
	if err != nil {
		return nil, err
	}
	return &API{db}, nil
}
