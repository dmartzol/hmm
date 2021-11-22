package postgres

import (
	"log"

	"github.com/dmartzol/hmm/internal/hmm"
)

func (db *DB) PopulateAccount(a *hmm.Account) *hmm.Account {
	var err error
	if a.Roles == nil {
		a.Roles, err = db.RolesForAccount(a.ID)
		if err != nil {
			log.Printf("PopulateAccount - RolesForAccount: %+v", err)
		}
	}
	return a
}

func (db *DB) PopulateAccounts(accs hmm.Accounts) hmm.Accounts {
	for _, a := range accs {
		db.PopulateAccount(a)
	}
	return accs
}
