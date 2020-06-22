package models

type Role struct {
	Row
	Name string
}

func (r Role) View(options map[string]bool) RoleView {
	view := RoleView{
		Name: r.Name,
	}
	return view
}

type RoleView struct {
	Name string
}

type CreateRoleReq struct {
	Name string
}
