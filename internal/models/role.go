package models

type Role struct {
	Row
	Name           string
	PermissionsBit RolePermission `db:"permission_bit"`
}

type Roles []*Role

func (r Role) View(options map[string]bool) RoleView {
	roleView := RoleView{
		Name: r.Name,
	}
	if options["permissions"] {
		for i := 1; i <= int(LastPermission); i *= 2 {
			if r.HasPermission(RolePermission(i)) {
				roleView.Permissions = append(roleView.Permissions, RolePermission(i).String())
			}
		}
	}
	return roleView
}

func (rs Roles) View(options map[string]bool) []RoleView {
	var views []RoleView
	for _, r := range rs {
		views = append(views, r.View(options))
	}
	return views
}

// HasPermission reports whether a role has the given permission
func (r Role) HasPermission(permission RolePermission) bool {
	if (r.PermissionsBit & permission) == permission {
		return true
	}
	return false
}

type RoleView struct {
	Name        string
	Permissions []string
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

type AddAccountRoleReq struct {
	AccountID, RoleID int64
}
