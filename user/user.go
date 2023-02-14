package user

var roles []string

func AddRole(role string) {
	for _, value := range roles {
		if value == role {
			roles = append(roles, role)

		}
	}
}

func manage(roles *role.Role) {

	roles.AddPermission("admin", "getPapers")
}
