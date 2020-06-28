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

func (r RolePermission) String() string {
	switch r {
	case PermissionAccountsView:
		return "view-accounts"
	case PermissionAccountsEdit:
		return "edit-accounts"
	default:
		return ""
	}
}

func (rp RolePermission) Int() int {
	return int(rp)
}
