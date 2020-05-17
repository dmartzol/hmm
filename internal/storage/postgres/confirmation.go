package postgres

import "github.com/dmartzol/hackerspace/internal/models"

func (db *DB) CreateConfirmation(accountID int64, t models.ConfirmationType) (*models.Confirmation, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var ecc models.Confirmation
	sqlStatement := `update confirmations set expire_time = current_timestamp where confirm_time is null and type = $1 and account_id = $2 returning *`
	err = tx.Select(&ecc, sqlStatement, t, accountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var cc models.Confirmation
	sqlStatement = `insert into confirmations (type, account_id) values ($1, $2) returning *`
	err = tx.Get(&cc, sqlStatement, t, accountID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &cc, tx.Commit()
}

func (db *DB) PendingConfirmationByKey(key string) (*models.Confirmation, error) {
	var c models.Confirmation
	sqlStatement := `select * from confirmations c where c."key" = $1 and c.confirm_time is null and c.expire_time > current_timestamp limit 1`
	err := db.Get(&c, sqlStatement, key)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *DB) FailedConfirmationIncrease(id int64) (*models.Confirmation, error) {
	var c models.Confirmation
	sqlStatement := `update confirmations set failed_confirmations_count = failed_confirmations_count + 1 where id = $1 returning *`
	err := db.Get(&c, sqlStatement, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *DB) Confirm(id int64) (*models.Confirmation, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	var c models.Confirmation
	sqlStatement := `update confirmations set confirm_time = current_timestamp where id = $1 and confirm_time is null returning *`
	err = tx.Get(&c, sqlStatement, id)
	if err != nil {
		return nil, err
	}
	var a models.Account
	sqlStatement = `update accounts set confirmed_email = true where id = $1 returning *`
	err = tx.Get(&a, sqlStatement, c.AccountID)
	if err != nil {
		return nil, err
	}
	return &c, tx.Commit()
}
