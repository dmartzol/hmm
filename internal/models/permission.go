package models

type RolePermission int

const (
	PermissionAccountsView RolePermission = 1 << iota
	PermissionAccountsEdit
	PermissionAccountsDeactivate
	PermissionAuthorizationsView
	PermissionAuthorizationsCreate
	PermissionAuthorizationsEdit
	PermissionAuthorizationsDelete
	PermissionRolesView
	PermissionRolesCreate
	PermissionRolesEdit
	PermissionRolesDelete
	LastPermission
)

func StringToRolePermission(s string) RolePermission {
	for i := 1; i <= int(LastPermission); i *= 2 {
		if RolePermission(i).String() == s {
			return RolePermission(i)
		}
	}
	return RolePermission(0)
}

func (r RolePermission) String() string {
	switch r {
	case PermissionAccountsView:
		return "view-accounts"
	case PermissionAccountsEdit:
		return "edit-accounts"
	case PermissionAccountsDeactivate:
		return "deactivate-accounts"
	case PermissionAuthorizationsView:
		return "view-authorizations"
	case PermissionAuthorizationsCreate:
		return "create-authorizations"
	case PermissionAuthorizationsEdit:
		return "edit-authorizations"
	case PermissionAuthorizationsDelete:
		return "delete-authorizations"
	case PermissionRolesView:
		return "view-roles"
	case PermissionRolesCreate:
		return "create-roles"
	case PermissionRolesEdit:
		return "edit-roles"
	case PermissionRolesDelete:
		return "delete-roles"
	default:
		return "unnamed-permission"
	}
}

func (rp RolePermission) Int() int {
	return int(rp)
}
