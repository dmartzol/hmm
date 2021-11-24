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

func (ss SessionService) SessionFromToken(token string) (*hmm.Session, error) {
	session, ok := ss.MemCache.SessionFromToken(token)
	if ok {
		return session, nil
	}
	session, err := ss.DB.SessionFromToken(token)
	if err != nil {
		return nil, err
	}
	ss.MemCache.AddSession(session)
	return session, nil
}

func (ss SessionService) ExpireSession(token string) (*hmm.Session, error) {
	session, err := ss.DB.ExpireSession(token)
	if err != nil {
		return nil, err
	}
	ss.MemCache.DeleteSession(token)
	return session, nil
}
