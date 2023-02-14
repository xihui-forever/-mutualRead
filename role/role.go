package role

type role struct {
	roleType string
	roleMap  map[string]bool
}

type Role struct {
	role        []*role
	permissions []string
}

var roleManage Role

func getRole(roleType string) *role {
	for _, value := range roleManage.role {
		if value.roleType == roleType {
			return value
		}
	}
	return nil
}

func AddRole(roleType string, permissions []string) error {
	ok := getRole(roleType)
	if ok != nil {
		return ErrRoleAlreadyExists
	}
	m := make(map[string]bool)
	for _, permission := range roleManage.permissions {
		m[permission] = false
	}
	for _, value := range permissions {
		_, ok := m[value]
		if !ok {
			return ErrPermissionNotExists
		}
		m[value] = true
	}
	roleManage.role = append(roleManage.role, &role{
		roleType: roleType,
		roleMap:  m,
	})
	return nil
}

func AddPermission(permission string) error {
	for _, value := range roleManage.permissions {
		if value == permission {
			return ErrPermisssionExists
		}
	}
	roleManage.permissions = append(roleManage.permissions, permission)
	return nil
}

func RemovePermission(permission string) error {
	for key, value := range roleManage.permissions {
		if value == permission {
			roleManage.permissions = append(roleManage.permissions[:key], roleManage.permissions[key+1:]...)
			return nil
		}
	}
	return ErrPermisssionExists
}

func AddRolePermission(roleType string, permission string) error {
	role := getRole(roleType)
	if role == nil {
		return ErrRoleNotExists
	}
	_, ok := role.roleMap[permission]
	if !ok {
		return ErrPermissionNotExists
	}
	role.roleMap[permission] = true
	return nil
}

func RemoveRolePermission(roleType string, permission string) error {
	role := getRole(roleType)
	if role == nil {
		return ErrRoleNotExists
	}
	_, ok := role.roleMap[permission]
	if !ok {
		return ErrPermissionNotExists
	}
	role.roleMap[permission] = false
	return nil
}
