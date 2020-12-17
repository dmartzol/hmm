package models

import "strings"

type Role struct {
	Row
	Name        string
	Permissions []string
}

type Roles []*Role

// View returns a role view
func (r Role) View(options map[string]bool) RoleView {
	roleView := RoleView{
		Name: r.Name,
	}
	roleView.Permissions = r.Permissions
	return roleView
}

// View returns role views
func (rs Roles) View(options map[string]bool) []RoleView {
	var views []RoleView
	for _, r := range rs {
		views = append(views, r.View(options))
	}
	return views
}

// HasPermission reports whether a role has the given permission
func (r Role) HasPermission(permission string) bool {
	for _, p := range r.Permissions {
		if strings.EqualFold(p, permission) {
			return true
		}
	}
	return false
}

type RoleView struct {
	Name          string
	Permissions   []string
	PermissionBit int
}

type AccountRole struct {
	Row
	AccountID int64 `db:"account_id"`
	RoleID    int64 `db:"role_id"`
}

func (ar AccountRole) View(options map[string]bool) AccountRoleView {
	view := AccountRoleView{}
	return view
}

type AccountRoleView struct {
	Account AccountView
	Role    RoleView
}

type CreateRoleReq struct {
	Name string
}

type EditRoleReq struct {
	Name        *string
	Permissions []string
}

type AddAccountRoleReq struct {
	AccountID, RoleID int64
}
