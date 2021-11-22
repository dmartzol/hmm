package handler

import "github.com/dmartzol/hmm/internal/storage/postgres"

const (
	apiVersionNumber = "0.0.1"
	hmmmCookieName   = "Hmm-Cookie"
	idQueryParameter = "id"
)

// API represents something
type Handler struct {
	db *postgres.DB
}

func New(db *postgres.DB) (*Handler, error) {
	return &Handler{db}, nil
}
