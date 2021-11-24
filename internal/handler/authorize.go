package handler

import (
	"fmt"
	"log"

	"github.com/dmartzol/hmm/internal/hmm"
)

func (h Handler) AuthorizeAccount(accountID int64, permission hmm.RolePermission) error {
	roles, err := h.RoleService.RolesForAccount(accountID)
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
