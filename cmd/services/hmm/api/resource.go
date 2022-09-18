package api

import (
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/internal/logger"
	"github.com/dmartzol/hmm/internal/storage"
	"github.com/dmartzol/hmm/internal/storage/postgres"
)

type Resources struct {
	AccountService      hmm.AccountService
	SessionService      hmm.SessionService
	ConfirmationService hmm.ConfirmationService
	RoleService         hmm.RoleService
	Logger              Logger
}

func newResources(structuredLogging bool, db *postgres.DB) *Resources {
	return &Resources{
		AccountService:      storage.NewAccountService(db),
		SessionService:      storage.NewSessionService(db),
		ConfirmationService: storage.NewConfirmationService(db),
		RoleService:         storage.NewRoleService(db),
		Logger:              logger.New(structuredLogging),
	}
}
