package models

type Permission string

const (
	PermissionAccountsView         string = "view-accounts"
	PermissionAccountsEdit                = "edit-accounts"
	PermissionAccountsDeactivate          = "deactivate-accounts"
	PermissionAuthorizationsView          = "view-authorizations"
	PermissionAuthorizationsCreate        = "create-authorizations"
	PermissionAuthorizationsEdit          = "edit-authorizations"
	PermissionAuthorizationsDelete        = "delete-authorizations"
	PermissionRolesView                   = "view-roles"
	PermissionRolesCreate                 = "create-roles"
	PermissionRolesEdit                   = "edit-roles"
	PermissionRolesDelete                 = "delete-roles"
)
