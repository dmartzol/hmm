package dao

import (
	"github.com/dmartzol/hmm/internal/dao/postgres"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/jmoiron/sqlx"
)

type ConfirmationService struct {
	DB *postgres.DB
}

func NewConfirmationService(db *sqlx.DB) *ConfirmationService {
	cs := ConfirmationService{
		DB: &postgres.DB{DB: db},
	}
	return &cs
}

func (cs ConfirmationService) PendingConfirmationByKey(key string) (*hmm.Confirmation, error) {
	conf, err := cs.DB.PendingConfirmationByKey(key)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (cs ConfirmationService) Confirm(id int64) (*hmm.Confirmation, error) {
	conf, err := cs.DB.Confirm(id)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (cs ConfirmationService) FailedConfirmationIncrease(id int64) (*hmm.Confirmation, error) {
	conf, err := cs.DB.FailedConfirmationIncrease(id)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
