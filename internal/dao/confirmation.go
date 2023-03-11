package dao

import (
	"github.com/dmartzol/hmm/internal/dao/memcache"
	"github.com/dmartzol/hmm/internal/dao/postgres"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/jmoiron/sqlx"
)

type ConfirmationService struct {
	DB       *postgres.DB
	MemCache *memcache.ConfirmationMemcache
}

func NewConfirmationService(db *sqlx.DB) *ConfirmationService {
	cs := ConfirmationService{
		DB:       &postgres.DB{DB: db},
		MemCache: memcache.NewConfirmationMemcache(),
	}
	return &cs
}

func (cs ConfirmationService) PendingConfirmationByKey(key string) (*hmm.Confirmation, error) {
	conf, ok := cs.MemCache.Confirmation(key)
	if ok {
		return conf, nil
	}
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
	cs.MemCache.Add(conf)
	return conf, nil
}

func (cs ConfirmationService) FailedConfirmationIncrease(id int64) (*hmm.Confirmation, error) {
	conf, err := cs.DB.FailedConfirmationIncrease(id)
	if err != nil {
		return nil, err
	}
	cs.MemCache.Add(conf)
	return conf, nil
}
