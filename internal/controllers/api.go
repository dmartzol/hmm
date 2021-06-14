package controllers

import "github.com/dmartzol/hmm/internal/storage/postgres"

const (
	apiVersionNumber = "0.0.1"
	hmmmCookieName   = "Hmm-Cookie"
	idQueryParameter = "id"
)

// API represents something
type API struct {
	db *postgres.DB
}

func NewAPI(db *postgres.DB) (*API, error) {
	return &API{db}, nil
}
