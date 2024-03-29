package postgres

import (
	"context"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"

	"github.com/dmartzol/hmm/internal/hmm"
)

// Account fetches an account by id
func (db *DB) Account(id int64) (*hmm.Account, error) {
	var a hmm.Account
	sqlStatement := `select * from accounts a where a.id = $1`
	err := db.Get(&a, sqlStatement, id)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Accounts returns all accounts in the db
func (db *DB) Accounts() (hmm.Accounts, error) {
	var accs []*hmm.Account
	sqlStatement := `select * from accounts a`
	err := db.Select(&accs, sqlStatement)
	if err != nil {
		return nil, err
	}
	return accs, nil
}

// AccountWithCredentials returns an account if the email and password provided match an (email,password) pair in the db
func (db *DB) AccountWithCredentials(email, allegedPassword string) (*hmm.Account, error) {
	var a hmm.Account
	sqlStatement := `select * from accounts a where a.email = $1 and a.passhash = crypt($2, a.passhash)`
	err := db.Get(&a, sqlStatement, email, allegedPassword)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CreateAccount creates a new account in the db and a confirmation code for the new registered email
func (db *DB) CreateAccount(ctx context.Context, a *hmm.Account, password, confirmationCode string) (*hmm.Account, *hmm.Confirmation, error) {
	_, span := otel.Tracer("db").Start(ctx, "DB.CreateAccount")
	defer span.End()

	tx, err := db.Beginx()
	if err != nil {
		return nil, nil, err
	}

	var newAccount hmm.Account
	sqlStatement := `
		INSERT INTO accounts (
			first_name,
			last_name,
			dob,
			gender,
			phone_number,
			email,
			passhash
		) values ($1, $2, $3, $4, $5, $6, crypt($7, gen_salt('bf', 8))) returning *
	`
	err = tx.Get(
		&newAccount,
		sqlStatement,
		a.FirstName,
		a.LastName,
		a.DOB,
		a.Gender,
		a.PhoneNumber,
		a.Email,
		password,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}

	var confirmation hmm.Confirmation
	sqlStatement = `
		INSERT INTO confirmations (
			type,
			account_id,
			key,
			confirmation_target
		) values ($1, $2, $3, $4) returning *
	`
	err = tx.Get(&confirmation, sqlStatement, hmm.ConfirmationTypeEmail, newAccount.ID, confirmationCode, newAccount.Email)
	if err != nil {
		_ = tx.Rollback()
		return nil, nil, err
	}

	return &newAccount, &confirmation, tx.Commit()
}
