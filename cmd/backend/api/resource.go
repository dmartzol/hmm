package api

import (
	storage "github.com/dmartzol/hmm/internal/dao"
	"github.com/dmartzol/hmm/internal/dao/postgres"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/internal/logger"
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
