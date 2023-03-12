package dao

import (
	"github.com/dmartzol/hmm/internal/dao/postgres"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/jmoiron/sqlx"
)

type AccountService struct {
	DB *postgres.DB
}

func NewAccountService(db *sqlx.DB) *AccountService {
	as := AccountService{
		DB: &postgres.DB{DB: db},
	}
	return &as
}

func (as AccountService) Account(id int64) (*hmm.Account, error) {
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
	return newAccount, confirmation, nil
}

func (as AccountService) PopulateAccount(account *hmm.Account) *hmm.Account {
	return as.DB.PopulateAccount(account)
}

func (as AccountService) PopulateAccounts(accounts hmm.Accounts) hmm.Accounts {
	return as.DB.PopulateAccounts(accounts)
}
