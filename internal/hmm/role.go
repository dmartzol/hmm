package hmm

type Roles []*Role

type Role struct {
	Row
	Name           string
	PermissionsBit RolePermission `db:"permission_bit"`

	// Synthetic field
	Permissions []string
}

type RoleService interface {
}

// Populate populates synthetic fields for the role structure
func (r *Role) Populate() *Role {
	for i := 1; i <= int(LastPermission); i *= 2 {
		if r.HasPermission(RolePermission(i)) {
			r.Permissions = append(r.Permissions, RolePermission(i).String())
		}
	}
	return r
}

// Populate populates synthetic fields for role structures
func (rs Roles) Populate() Roles {
	for _, r := range rs {
		r.Populate()
	}
	return rs
}

// HasPermission reports whether a role has the given permission
func (r Role) HasPermission(permission RolePermission) bool {
	if (r.PermissionsBit & permission) == permission {
		return true
	}
	return false
}

type AccountRole struct {
	Row
	AccountID int64 `db:"account_id"`
	RoleID    int64 `db:"role_id"`
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
