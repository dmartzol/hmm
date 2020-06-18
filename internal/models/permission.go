package models

type RolePermission string

const (
	PermissionAccountsView        RolePermission = "accounts-view"
	PermissionAccountsEdit                       = "accounts-edit"
	PermissionAccountsDeactivate                 = "accounts-deactivate"
	PermissionAuthorizationCreate                = "authorizations-create"
)
