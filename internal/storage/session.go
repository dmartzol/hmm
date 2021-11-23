package storage

import (
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/internal/storage/memcache"
	"github.com/dmartzol/hmm/internal/storage/postgres"
)

type SessionService struct {
	MemCache *memcache.SessionMemcache
	DB       *postgres.DB
}

func NewSessionService(db *postgres.DB) *SessionService {
	ss := SessionService{
		DB:       db,
		MemCache: memcache.NewSessionMemcache(),
	}
	return &ss
}

func (ss SessionService) Create(email, password string) (*hmm.Session, error) {
	session, err := ss.DB.CreateSession(email, password)
	if err != nil {
		return nil, err
	}
	ss.MemCache.AddSession(session)
	return session, nil
}
