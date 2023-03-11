package api

import (
	improvedLogger "github.com/dmartzol/go-sdk/logger"
	storage "github.com/dmartzol/hmm/internal/dao"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/jmoiron/sqlx"
)

type Resources struct {
	AccountService      hmm.AccountService
	SessionService      hmm.SessionService
	ConfirmationService hmm.ConfirmationService
	RoleService         hmm.RoleService
	Logger              improvedLogger.Logger
}

func newResources(db *sqlx.DB, loggerr improvedLogger.Logger) *Resources {
	return &Resources{
		AccountService:      storage.NewAccountService(db),
		SessionService:      storage.NewSessionService(db),
		ConfirmationService: storage.NewConfirmationService(db),
		RoleService:         storage.NewRoleService(db),
		Logger:              loggerr,
	}
}
