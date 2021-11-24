package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dmartzol/hmm/internal/api"
	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/dmartzol/hmm/pkg/httpresponse"
	"github.com/gorilla/mux"
)

func (h Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req hmm.CreateRoleReq
	err := httpresponse.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	role, err := h.RoleService.Create(req.Name)
	if err != nil {
		h.Logger.Errorf("unable to create role: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	httpresponse.RespondJSON(w, api.RoleView(role, nil))
}

func (h Handler) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.RoleService.Roles()
	if err != nil {
		log.Printf("GetRoles Roles ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, api.RolesView(roles, nil))
}

func validateEditRole(req hmm.EditRoleReq, targetRole *hmm.Role) error {
	if req.Name == nil && len(req.Permissions) == 0 {
		return fmt.Errorf("No updates found")
	}
	return nil
}

func (h Handler) EditRole(w http.ResponseWriter, r *http.Request) {
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
	err = h.AuthorizeAccount(requesterID, hmm.PermissionRolesEdit)
	if err != nil {
		log.Printf("WARNING: account %d requested to edit role %d", requesterID, roleID)
		log.Printf("EditRole AuthorizeAccount ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	var req hmm.EditRoleReq
	err = httpresponse.Unmarshal(r, &req)
	if err != nil {
		log.Printf("CreateRole Unmarshal ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	role, err := h.RoleService.Role(roleID)
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

	newBit := 0
	for _, p := range req.Permissions {
		newBit = newBit | hmm.StringToRolePermission(p).Int()
	}
	if role.PermissionsBit.Int() == newBit {
		log.Printf("EditRole ERROR: role already has those permissions")
		httpresponse.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	updatedRole, err := h.RoleService.Update(role.ID, newBit)
	if err != nil {
		log.Printf("EditRole UpdateRole ERROR: %+v", err)
		httpresponse.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	httpresponse.RespondJSON(w, api.RoleView(updatedRole, nil))
}
