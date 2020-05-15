package postgres

import (
	"time"

	"github.com/dmartzol/hackerspace/internal/models"
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

// Account fetches an account by id
func (db *DB) Account(id int64) (*models.Account, error) {
	var a models.Account
	sqlStatement := `select * from accounts a where a.id = $1`
	err := db.Get(&a, sqlStatement, id)
	if err != nil {
		return nil, err
	}
	return &a, nil
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

// AccountWithCredentials returns an account if the email and password provided match an (email,password) pair in the db
func (db *DB) AccountWithCredentials(email, allegedPassword string) (*models.Account, error) {
	var a models.Account
	sqlStatement := `select * from accounts a where a.email = $1 and a.passhash = crypt($2, a.passhash)`
	err := db.Get(&a, sqlStatement, email, allegedPassword)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CreateAccount creates a new account in the db and a confirmation code for the new registered email
func (db *DB) CreateAccount(first, last, email, password string, dob time.Time, gender, phone *string) (*models.Account, *models.ConfirmationCode, error) {
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
	var cc models.ConfirmationCode
	sqlStatement = `insert into confirmation_codes (type, account_id) values ($1, $2) returning *`
	err = tx.Get(&cc, sqlStatement, models.ConfirmationCodeTypeEmail, a.ID)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	return &a, &cc, tx.Commit()
}

func (db *DB) CreateConfirmationCode(accountID int64, t models.ConfirmationCodeType) (*models.ConfirmationCode, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var ecc models.ConfirmationCode
	sqlStatement := `update confirmation_codes set expire_time = current_timestamp where confirm_time is null and type = $1 and account_id = $2 returning *`
	err = tx.Select(&ecc, sqlStatement, t, accountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var cc models.ConfirmationCode
	sqlStatement = `insert into confirmation_codes (type, account_id) values ($1, $2) returning *`
	err = tx.Get(&cc, sqlStatement, t, accountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &cc, tx.Commit()
}
