package postgres

import "github.com/dmartzol/hmm/internal/hmm"

// CreateRole creates a new role with the given name
func (db *DB) CreateRole(name string) (*hmm.Role, error) {
	var r hmm.Role
	sqlStatement := `insert into roles (name, permission_bit) values ($1, 0) returning *`
	err := db.Get(&r, sqlStatement, name)
	if err != nil {
		return nil, err
	}
	return r.Populate(), nil
}

func (db *DB) Role(id int64) (*hmm.Role, error) {
	var r hmm.Role
	sqlSelect := `select * from roles where id = $1`
	err := db.Get(&r, sqlSelect, id)
	if err != nil {
		return nil, err
	}
	return r.Populate(), nil
}

func (db *DB) AddRoleToAccount(roleID, accountID int64) (*hmm.AccountRole, error) {
	var r hmm.AccountRole
	sqlInsert := `insert into account_roles (role_id, account_id) values ($1, $2) returning *`
	err := db.Get(&r, sqlInsert, roleID, accountID)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// RolesForAccount fetches all roles for the given account
func (db *DB) RolesForAccount(accountID int64) (hmm.Roles, error) {
	var rs hmm.Roles
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
func (db *DB) Roles() (hmm.Roles, error) {
	var rs hmm.Roles
	sqlStatement := `select * from roles`
	err := db.Select(&rs, sqlStatement)
	if err != nil {
		return nil, err
	}
	return rs.Populate(), nil
}

func (db *DB) UpdateRole(roleID int64, permissionBit int) (*hmm.Role, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var r hmm.Role
	sqlStatement := `update roles set permission_bit = $1 where id = $2 returning *`
	err = tx.Get(&r, sqlStatement, permissionBit, roleID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return r.Populate(), tx.Commit()
}
