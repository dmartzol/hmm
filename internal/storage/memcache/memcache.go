package memcache

import (
	"errors"

	"github.com/dmartzol/hmm/internal/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

type AccountMemcache map[int64]*domain.Account

func NewAccountMemcache() *AccountMemcache {
	m := make(AccountMemcache)
	return &m
}

func (a AccountMemcache) Account(id int64) (*domain.Account, bool) {
	acc, ok := a[id]
	if !ok {
		return nil, false
	}
	return acc, true
}
