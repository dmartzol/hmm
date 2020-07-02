package postgres

import "github.com/dmartzol/hmmm/internal/models"

// RoleExists returns true if already exists a role with the provided name in the db
func (db *DB) RoleExists(name string) (bool, error) {
	var exists bool
	sqlStatement := `select exists(select 1 from roles r where r.name = $1)`
	err := db.Get(&exists, sqlStatement, name)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateRole creates a new role with the given name
func (db *DB) CreateRole(name string) (*models.Role, error) {
	var r models.Role
	sqlStatement := `insert into roles (name, permission_bit) values ($1, 0) returning *`
	err := db.Get(&r, sqlStatement, name)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (db *DB) Role(roleID int64) (*models.Role, error) {
	var r models.Role
	sqlStatement := `select * from roles where id = $1`
	err := db.Get(&r, sqlStatement, roleID)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (db *DB) AddAccountRole(roleID, accountID int64) (*models.AccountRole, error) {
	var r models.AccountRole
	sqlStatement := `insert into account_roles (role_id, account_id) values ($1, $2) returning *`
	err := db.Get(&r, sqlStatement, roleID, accountID)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// RolesForAccount fetches all roles for the given account
func (db *DB) RolesForAccount(accountID int64) (models.Roles, error) {
	var rs []*models.Role
	sqlStatement := `select r.* from roles r 
	inner join account_roles ar on ar.role_id = r.id
	where
	ar.account_id = $1`
	err := db.Select(&rs, sqlStatement, accountID)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// Roles fetches all roles in the database
func (db *DB) Roles() (models.Roles, error) {
	var rs models.Roles
	sqlStatement := `select * from roles`
	err := db.Select(&rs, sqlStatement)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
