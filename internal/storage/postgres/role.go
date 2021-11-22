package postgres

import "github.com/dmartzol/hmm/internal/domain"

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
func (db *DB) CreateRole(name string) (*domain.Role, error) {
	var r domain.Role
	sqlStatement := `insert into roles (name, permission_bit) values ($1, 0) returning *`
	err := db.Get(&r, sqlStatement, name)
	if err != nil {
		return nil, err
	}
	return r.Populate(), nil
}

func (db *DB) Role(roleID int64) (*domain.Role, error) {
	var r domain.Role
	sqlStatement := `select * from roles where id = $1`
	err := db.Get(&r, sqlStatement, roleID)
	if err != nil {
		return nil, err
	}
	return r.Populate(), nil
}

func (db *DB) AddAccountRole(roleID, accountID int64) (*domain.AccountRole, error) {
	var r domain.AccountRole
	sqlStatement := `insert into account_roles (role_id, account_id) values ($1, $2) returning *`
	err := db.Get(&r, sqlStatement, roleID, accountID)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// RolesForAccount fetches all roles for the given account
func (db *DB) RolesForAccount(accountID int64) (domain.Roles, error) {
	var rs domain.Roles
	sqlStatement := `select r.* from roles r 
	inner join account_roles ar on ar.role_id = r.id
	where
	ar.account_id = $1`
	err := db.Select(&rs, sqlStatement, accountID)
	if err != nil {
		return nil, err
	}
	return rs.Populate(), nil
}

// Roles fetches all roles in the database
func (db *DB) Roles() (domain.Roles, error) {
	var rs domain.Roles
	sqlStatement := `select * from roles`
	err := db.Select(&rs, sqlStatement)
	if err != nil {
		return nil, err
	}
	return rs.Populate(), nil
}

func (db *DB) UpdateRole(roleID int64, permissionBit int) (*domain.Role, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var r domain.Role
	sqlStatement := `update roles set permission_bit = $1 where id = $2 returning *`
	err = tx.Get(&r, sqlStatement, permissionBit, roleID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return r.Populate(), tx.Commit()
}
