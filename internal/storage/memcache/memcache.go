package memcache

import (
	"errors"

	"github.com/dmartzol/hmm/internal/models"
)

var (
	ErrNotFound = errors.New("not found")
)

type AccountMemcache map[int64]*models.Account

func NewAccountMemcache() *AccountMemcache {
	m := make(AccountMemcache)
	return &m
}

func (a AccountMemcache) Account(id int64) (*models.Account, bool) {
	acc, ok := a[id]
	if !ok {
		return nil, false
	}
	return acc, true
}
