package storage

import (
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/internal/storage/memcache"
	"github.com/dmartzol/hmm/internal/storage/postgres"
)

type RoleService struct {
	MemCache *memcache.RoleMemcache
	DB       *postgres.DB
}

func NewRoleService(db *postgres.DB) *RoleService {
	rs := RoleService{
		DB:       db,
		MemCache: memcache.NewRoleMemcache(),
	}
	return &rs
}

func (rs RoleService) Role(id int64) (*hmm.Role, error) {
	role, ok := rs.MemCache.Role(id)
	if ok {
		return role, nil
	}
	role, err := rs.DB.Role(id)
	if err != nil {
		return nil, err
	}
	return role, nil
}
