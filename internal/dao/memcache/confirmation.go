package memcache

import "github.com/dmartzol/hmm/internal/hmm"

type ConfirmationMemcache map[string]*hmm.Confirmation

func NewConfirmationMemcache() *ConfirmationMemcache {
	m := make(ConfirmationMemcache)
	return &m
}

func (m ConfirmationMemcache) Confirmation(key string) (*hmm.Confirmation, bool) {
	conf, ok := m[key]
	if !ok {
		return nil, false
	}
	return conf, true
}

func (m ConfirmationMemcache) Add(conf *hmm.Confirmation) {
	m[conf.Key] = conf
}
