package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dmartzol/hmm/internal/hmm"
	"github.com/gorilla/mux"
)

type AccountRole struct {
	Account Account
	Role    Role
}

func AccountRoleView(ar *hmm.AccountRole, options map[string]bool) AccountRole {
	view := AccountRole{}
	return view
}

type Role struct {
	Name          string
	Permissions   []string
	PermissionBit int
}

func RolesView(rs hmm.Roles, options map[string]bool) []Role {
	var views []Role
	for _, r := range rs {
		views = append(views, RoleView(r, options))
	}
	return views
}

func RoleView(r *hmm.Role, options map[string]bool) Role {
	roleView := Role{
		Name:          r.Name,
		PermissionBit: r.PermissionsBit.Int(),
	}
	if len(r.Permissions) == 0 {
		r.Populate()
	}
	roleView.Permissions = r.Permissions
	return roleView
}

func (h API) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req hmm.CreateRoleReq
	err := h.Unmarshal(r, &req)
	if err != nil {
		h.Logger.Errorf("unable to unmarshal: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	role, err := h.RoleService.Create(req.Name)
	if err != nil {
		h.Logger.Errorf("unable to create role: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	h.RespondJSON(w, RoleView(role, nil))
}

func (h API) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.RoleService.Roles()
	if err != nil {
		log.Printf("GetRoles Roles ERROR: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	h.RespondJSON(w, RolesView(roles, nil))
}

func validateEditRole(req hmm.EditRoleReq, targetRole *hmm.Role) error {
	if req.Name == nil && len(req.Permissions) == 0 {
		return fmt.Errorf("No updates found")
	}
	return nil
}

func (h API) EditRole(w http.ResponseWriter, r *http.Request) {
	// parsing parameters
	params := mux.Vars(r)
	idString, ok := params[idQueryParameter]
	if !ok {
		errMsg := fmt.Sprintf("parameter '%s' not found", idQueryParameter)
		h.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}
	roleID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("wrong parameter '%s'", idString)
		h.RespondJSONError(w, errMsg, http.StatusBadRequest)
		return
	}
	// checking permissions
	ctx := r.Context()
	requesterID := ctx.Value(contextRequesterAccountIDKey).(int64)
	err = h.AuthorizeAccount(requesterID, hmm.PermissionRolesEdit)
	if err != nil {
		log.Printf("WARNING: account %d requested to edit role %d", requesterID, roleID)
		log.Printf("EditRole AuthorizeAccount ERROR: %+v", err)
		h.RespondJSONError(w, "", http.StatusUnauthorized)
		return
	}

	var req hmm.EditRoleReq
	err = h.Unmarshal(r, &req)
	if err != nil {
		log.Printf("CreateRole Unmarshal ERROR: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}

	role, err := h.RoleService.Role(roleID)
	if err != nil {
		log.Printf("EditRole Role ERROR: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	err = validateEditRole(req, role)
	if err != nil {
		log.Printf("EditRole validateEditRole ERROR: %+v", err)
		h.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}

	newBit := 0
	for _, p := range req.Permissions {
		newBit = newBit | hmm.StringToRolePermission(p).Int()
	}
	if role.PermissionsBit.Int() == newBit {
		log.Printf("EditRole ERROR: role already has those permissions")
		h.RespondJSONError(w, "", http.StatusBadRequest)
		return
	}
	updatedRole, err := h.RoleService.Update(role.ID, newBit)
	if err != nil {
		log.Printf("EditRole UpdateRole ERROR: %+v", err)
		h.RespondJSONError(w, "", http.StatusInternalServerError)
		return
	}
	h.RespondJSON(w, RoleView(updatedRole, nil))
}
