package api

import (
	"github.com/dmartzol/go-sdk/logger"
	"github.com/dmartzol/hmm/internal/dao"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/jmoiron/sqlx"
)

type Resources struct {
	AccountService      hmm.AccountService
	SessionService      hmm.SessionService
	ConfirmationService hmm.ConfirmationService
	RoleService         hmm.RoleService
	Logger              logger.Logger
}

func newResources(db *sqlx.DB, l logger.Logger) *Resources {
	return &Resources{
		AccountService:      dao.NewAccountService(db),
		SessionService:      dao.NewSessionService(db),
		ConfirmationService: dao.NewConfirmationService(db),
		RoleService:         dao.NewRoleService(db),
		Logger:              l,
	}
}
