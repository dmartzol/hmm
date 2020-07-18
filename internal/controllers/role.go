package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dmartzol/hmm/internal/models"
	"github.com/dmartzol/hmm/pkg/httpresponse"
	"github.com/gorilla/mux"
)

type roleStorage interface {
	CreateRole(name string) (*models.Role, error)
	Roles() (models.Roles, error)
	RoleExists(name string) (bool, error)
	RolesForAccount(accountID int64) (models.Roles, error)
	AddAccountRole(roleID, accountID int64) (*models.AccountRole, error)
	Role(roleID int64) (*models.Role, error)
	UpdateRole(roleID int64, permissionBit int) (*models.Role, error)
}

func (api API) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("CreateRole Unmarshal ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	exists, err := api.RoleExists(req.Name)
	if err != nil {
		log.Printf("CreateRole RoleExists ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		log.Printf("ERROR: Role '%s' already exists", req.Name)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}
	role, err := api.storage.CreateRole(req.Name)
	if err != nil {
		log.Printf("CreateRole storage.CreateRole ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, role.View(nil))
}

func (api API) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := api.Roles()
	if err != nil {
		log.Printf("GetRoles Roles ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("Role: %+v", roles[0])
	log.Printf("Role: %d", roles[0].PermissionsBit)
	httpresponse.RespondJSON(w, roles.View(nil))
}

func validateEditRole(req models.EditRoleReq, targetRole *models.Role) error {
	if req.PermissionsBit == nil && req.Name == nil && len(req.Permissions) == 0 {
		return fmt.Errorf("No updates found")
	}
	return nil
}

func (api API) EditRole(w http.ResponseWriter, r *http.Request) {
	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params["id"]
	if !ok {
		http.Error(w, "parameter 'id' not found", http.StatusBadRequest)
		return
	}
	roleID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("wrong parameter '%s'", idString), http.StatusBadRequest)
		return
	}
	// checking permissions
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	err = api.AuthorizeAccount(requesterID, models.PermissionRolesEdit)
	if err != nil {
		log.Printf("WARNING: account %d requested to edit role %d", requesterID, roleID)
		log.Printf("EditRole AuthorizeAccount ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var req models.EditRoleReq
	err = httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("CreateRole Unmarshal ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	role, err := api.Role(roleID)
	if err != nil {
		log.Printf("EditRole Role ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = validateEditRole(req, role)
	if err != nil {
		log.Printf("EditRole validateEditRole ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// TODO: pass permission by string in request instead of int
	newPermissionBit := models.RolePermission(1)
	for _, p := range req.Permissions {
		newPermissionBit = newPermissionBit | models.StringToRolePermission(p)
	}
	if role.HasPermission(newPermissionBit) {
		log.Printf("EditRole ERROR: role already has those permissions")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	updatedRole, err := api.UpdateRole(role.ID, newPermissionBit.Int())
	if err != nil {
		log.Printf("EditRole UpdateRole ERROR: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, updatedRole.View(nil))
}
