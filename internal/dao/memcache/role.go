package memcache

import "github.com/dmartzol/hmm/internal/hmm"

type RoleMemcache map[int64]*hmm.Role

func NewRoleMemcache() *RoleMemcache {
	m := make(RoleMemcache)
	return &m
}

func (m RoleMemcache) Role(id int64) (*hmm.Role, bool) {
	role, ok := m[id]
	if !ok {
		return nil, false
	}
	return role, true
}

func (m RoleMemcache) Add(role *hmm.Role) {
	m[role.ID] = role
}
