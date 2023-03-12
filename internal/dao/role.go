package dao

import (
	"github.com/dmartzol/hmm/internal/dao/postgres"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/jmoiron/sqlx"
)

type RoleService struct {
	DB *postgres.DB
}

func NewRoleService(db *sqlx.DB) *RoleService {
	rs := RoleService{
		DB: &postgres.DB{DB: db},
	}
	return &rs
}

func (rs RoleService) Role(id int64) (*hmm.Role, error) {
	role, err := rs.DB.Role(id)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (rs RoleService) Roles() (hmm.Roles, error) {
	roles, err := rs.DB.Roles()
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (rs RoleService) RolesForAccount(id int64) (hmm.Roles, error) {
	roles, err := rs.DB.RolesForAccount(id)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (rs RoleService) Create(name string) (*hmm.Role, error) {
	newRole, err := rs.DB.CreateRole(name)
	if err != nil {
		return nil, err
	}
	return newRole, nil
}

func (rs RoleService) Update(id int64, permissionBit int) (*hmm.Role, error) {
	role, err := rs.DB.UpdateRole(id, permissionBit)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (rs RoleService) AddRoleToAccount(accountID int64, roleID int64) (*hmm.AccountRole, error) {
	accountRole, err := rs.DB.AddRoleToAccount(accountID, roleID)
	if err != nil {
		return nil, err
	}
	return accountRole, nil
}
