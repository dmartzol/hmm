package postgres

import (
	"time"

	"github.com/dmartzol/hmm/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// AccountExists returns true if the provided email exists in the db
func (db *DB) AccountExists(email string) (bool, error) {
	var exists bool
	sqlStatement := `select exists(select 1 from accounts a where a.email = $1)`
	err := db.Get(&exists, sqlStatement, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DB) EmailExistsUsingTx(tx *sqlx.Tx, email string) (bool, error) {
	var exists bool
	sql := `select exists(select 1 from emails e where e.email = $1)`
	err := tx.Get(&exists, sql, email)
	return exists, err
}

// Account fetches an account by id
func (db *DB) Account(id int64) (*models.Account, error) {
	var a models.Account
	sqlStatement := `select * from accounts a where a.id = $1`
	err := db.Get(&a, sqlStatement, id)
	return &a, err
}

// AccountUsingTx fetches an account by id using a database tx
func (db *DB) AccountUsingTx(tx *sqlx.Tx, id int64) (*models.Account, error) {
	var a models.Account
	sqlStatement := `select * from accounts a where a.id = $1`
	err := tx.Get(&a, sqlStatement, id)
	return &a, err
}

// Accounts returns all accounts in the db
func (db *DB) Accounts() (models.Accounts, error) {
	var accs []*models.Account
	sqlStatement := `select * from accounts a`
	err := db.Select(&accs, sqlStatement)
	if err != nil {
		return nil, err
	}
	return accs, nil
}

// AccountFromCredentials returns an account if the email and password provided match an (email,password) pair in the db
func (db *DB) AccountFromCredentials(email, password string) (*models.Account, error) {
	tx, err := db.Beginx()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	a, err := db.AccountFromCredentialsUsingTx(tx, email, password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return a, tx.Commit()
}

// AccountFromCredentials returns an account if the email and password provided match an (email,password) pair in the db
func (db *DB) AccountFromCredentialsUsingTx(tx *sqlx.Tx, email, password string) (*models.Account, error) {
	var a models.Account
	sqlStatement := `select a.* from accounts a inner join emails e ON a.email_id = e.id where e.email = $1 and a.passhash = crypt($2, a.passhash)`
	err := tx.Get(&a, sqlStatement, email, password)
	return &a, err
}

// CreateAccount creates a new account in the db and a confirmation code for the new registered email
func (db *DB) CreateAccount(first, last, email, password, confirmationCode string, dob time.Time, gender, phone *string) (*models.Account, *models.Confirmation, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, nil, err
	}
	var a models.Account
	sqlStatement := `insert into accounts (first_name, last_name, dob, gender, phone_number, email, passhash) values ($1, $2, $3, $4, $5, $6, crypt($7, gen_salt('bf', 8))) returning *`
	err = tx.Get(&a, sqlStatement, first, last, dob, gender, phone, email, password)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	var cc models.Confirmation
	sqlStatement = `insert into confirmations (type, account_id, key, confirmation_target) values ($1, $2, $3, $4) returning *`
	err = tx.Get(&cc, sqlStatement, models.ConfirmationTypeEmail, a.ID, confirmationCode, email)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	return &a, &cc, tx.Commit()
}
