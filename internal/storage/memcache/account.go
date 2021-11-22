package memcache

import (
	"github.com/dmartzol/hmm/internal/hmm"
)

type AccountMemcache map[int64]*hmm.Account

func NewAccountMemcache() *AccountMemcache {
	m := make(AccountMemcache)
	return &m
}

func (a AccountMemcache) Account(id int64) (*hmm.Account, bool) {
	acc, ok := a[id]
	if !ok {
		return nil, false
	}
	return acc, true
}
