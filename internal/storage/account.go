package storage

import (
	"github.com/dmartzol/hmm/internal/hmm"
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

func (as AccountService) Account(id int64) (*hmm.Account, error) {
	account, ok := as.MemCache.Account(id)
	if ok {
		return account, nil
	}
	account, err := as.DB.Account(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (as AccountService) Accounts() (hmm.Accounts, error) {
	accs, err := as.DB.Accounts()
	if err != nil {
		return nil, err
	}
	return accs, nil
}

func (as AccountService) Create(account *hmm.Account, password, confirmationCode string) (*hmm.Account, *hmm.Confirmation, error) {
	newAccount, confirmation, err := as.DB.CreateAccount(account, password, confirmationCode)
	if err != nil {
		return nil, nil, err
	}
	as.MemCache.AddAccount(newAccount)
	return newAccount, confirmation, nil
}

func (as AccountService) PopulateAccount(account *hmm.Account) *hmm.Account {
	return as.DB.PopulateAccount(account)
}

func (as AccountService) PopulateAccounts(accounts hmm.Accounts) hmm.Accounts {
	return as.DB.PopulateAccounts(accounts)
}
