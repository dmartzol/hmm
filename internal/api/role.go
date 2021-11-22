package api

import "github.com/dmartzol/hmm/internal/hmm"

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
