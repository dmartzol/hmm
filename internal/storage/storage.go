package storage

import (
	"github.com/dmartzol/hmm/internal/models"
	"github.com/dmartzol/hmm/internal/storage/memcache"
	"github.com/dmartzol/hmm/internal/storage/postgres"
)

type AccountService struct {
	MemCache *memcache.AccountMemcache
	DB       *postgres.DB
}

func NewAccountService(db *postgres.DB) *AccountService {
	as := AccountService{
		DB:       db,
		MemCache: memcache.NewAccountMemcache(),
	}
	return &as
}

func (a AccountService) Account(id int64) (*models.Account, error) {
	account, ok := a.MemCache.Account(id)
	if ok {
		return account, nil
	}
	account, err := a.DB.Account(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a AccountService) Accounts() (models.Accounts, error) {
	panic("not implemented")
}

func (a AccountService) Create(req models.RegisterRequest) (*models.Account, error) {
	panic("not implemented")
}
