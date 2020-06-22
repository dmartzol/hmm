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
	sqlStatement := `insert into roles (name) values ($1) returning *`
	err := db.Get(&r, sqlStatement, name)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
