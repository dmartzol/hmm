package models

type RolePermission int

const (
	PermissionAccountsView        RolePermission = 1
	PermissionAccountsEdit                       = 2
	PermissionAccountsDeactivate                 = 4
	PermissionAuthorizationCreate                = 8
)
