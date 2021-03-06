package controllers

import (
	"fmt"
	"log"

	"github.com/dmartzol/hmm/internal/models"
)

func (api API) AuthorizeAccount(accountID int64, permission models.RolePermission) error {
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
