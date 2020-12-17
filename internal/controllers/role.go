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

func (api API) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("CreateRole Unmarshal ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	exists, err := api.db.RoleExists(req.Name)
	if err != nil {
		log.Printf("CreateRole RoleExists ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	// TODO: Use SQL tx for this
	if exists {
		log.Printf("ERROR: Role '%s' already exists", req.Name)
		httpresponse.RespondJSONError(w, "", http.StatusConflict)
		return
	}
	role, err := api.db.CreateRole(req.Name)
	if err != nil {
		log.Printf("CreateRole storage.CreateRole ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, role.View(nil))
}

func (api API) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := api.db.Roles()
	if err != nil {
		log.Printf("GetRoles Roles ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, roles.View(nil))
}

func validateEditRole(req models.EditRoleReq, targetRole *models.Role) error {
	if req.Name == nil && len(req.Permissions) == 0 {
		return fmt.Errorf("No updates found")
	}
	return nil
}

func (api API) EditRole(w http.ResponseWriter, r *http.Request) {
	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMsg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}
	roleID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		httpresponse.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}
	// checking permissions
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	err = api.AuthorizeAccount(requesterID, models.PermissionRolesEdit)
	if err != nil {
		log.Printf("WARNING: account %d requested to edit role %d", requesterID, roleID)
		log.Printf("EditRole AuthorizeAccount ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	var req models.EditRoleReq
	err = httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("CreateRole Unmarshal ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	role, err := api.db.Role(roleID)
	if err != nil {
		log.Printf("EditRole Role ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	err = validateEditRole(req, role)
	if err != nil {
		log.Printf("EditRole validateEditRole ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

	// newBit := 0
	// for _, p := range req.Permissions {
	// 	newBit = newBit | models.StringToRolePermission(p).Int()
	// }
	// TODO: update
	// if role.PermissionsBit.Int() == newBit {
	// 	log.Printf("EditRole ERROR: role already has those permissions")
	// 	httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
	// 	return
	// }
	updatedRole, err := api.db.UpdateRole(role.ID, req)
	if err != nil {
		log.Printf("EditRole UpdateRole ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, updatedRole.View(nil))
}
