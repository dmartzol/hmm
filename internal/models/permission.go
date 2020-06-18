package models

type RolePermission int

const (
	PermissionAccountsView RolePermission = iota
	PermissionAccountsEdit
	PermissionAccountsDeactivate
	PermissionAuthorizationAdd
)
