package role

import "errors"

var (
	ErrRolePermExists = errors.New("rolePerm already exists")
	ErrPermNotExist   = errors.New("permission not exists")
)
