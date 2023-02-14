package role

import "errors"

var (
	ErrRoleAlreadyExists   = errors.New("role already exists")
	ErrPermisssionExists   = errors.New("permisssion already exists")
	ErrPermissionNotExists = errors.New("permission not exists")
	ErrRoleNotExists       = errors.New("role not exists")
)
