package postgres

import (
	"log"

	"github.com/dmartzol/hmm/internal/models"
)

func (db *DB) PopulateAccount(a *models.Account) *models.Account {
	var err error
	if a.Roles == nil {
		a.Roles, err = db.RolesForAccount(a.ID)
		if err != nil {
			log.Printf("PopulateAccount - RolesForAccount: %+v", err)
		}
	}
	return a
}

func (db *DB) PopulateAccounts(accs models.Accounts) models.Accounts {
	for _, a := range accs {
		db.PopulateAccount(a)
	}
	return accs
}
