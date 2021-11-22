package postgres

import (
	"log"

	"github.com/dmartzol/hmm/internal/domain"
)

func (db *DB) PopulateAccount(a *domain.Account) *domain.Account {
	var err error
	if a.Roles == nil {
		a.Roles, err = db.RolesForAccount(a.ID)
		if err != nil {
			log.Printf("PopulateAccount - RolesForAccount: %+v", err)
		}
	}
	return a
}

func (db *DB) PopulateAccounts(accs domain.Accounts) domain.Accounts {
	for _, a := range accs {
		db.PopulateAccount(a)
	}
	return accs
}
