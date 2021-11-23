package handler

import (
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/internal/logger"
	"github.com/dmartzol/hmm/internal/storage"
	"github.com/dmartzol/hmm/internal/storage/postgres"
)

const (
	apiVersionNumber = "0.0.1"
	hmmmCookieName   = "Hmm-Cookie"
	idQueryParameter = "id"
)

// API represents something
type Handler struct {
	db             *postgres.DB
	AccountService hmm.AccountService
	SessionService hmm.SessionService
	Logger         Logger
}

func New(structuredLogging bool, db *postgres.DB) (*Handler, error) {
	handler := Handler{
		db:             db,
		AccountService: storage.NewAccountService(db),
		SessionService: storage.NewSessionService(db),
	}
	handler.Logger = logger.New(structuredLogging)
	return &handler, nil
}
