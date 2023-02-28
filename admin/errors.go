package admin

import "errors"

var (
	ErrAdminExist           = errors.New("admin exist")
	ErrAdminNotExist        = errors.New("admin not exist")
	ErrAdminRemoveFailed    = errors.New("admin remove failed")
	ErrPasswordWrong        = errors.New("password wrong")
	ErrPasswordChangeFailed = errors.New("password change failed")
)
