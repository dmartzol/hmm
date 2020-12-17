package controllers

import (
	"fmt"
	"log"
)

func (api API) AuthorizeAccount(accountID int64, permission string) error {
	roles, err := api.db.RolesForAccount(accountID)
	if err != nil {
		log.Printf("AuthorizeAccount Account ERROR: %+v", err)
		return err
	}
	for _, role := range roles {
		if role.HasPermission(permission) {
			return nil
		}
	}
	return fmt.Errorf("account %d not authorized to %s", accountID, permission)
}
