package mock

import "github.com/dmartzol/hmm/internal/hmm"

type AccountService struct {
	AccountFn func(string) (*hmm.Account, error)
}

func (as *AccountService) Account(id string) (*hmm.Account, error) {
	return as.AccountFn(id)
}
