package student

import "errors"

var (
	ErrStudentExist         = errors.New("student already exists")
	ErrStudentNotExist      = errors.New("student not exists")
	ErrPasswordWrong        = errors.New("password is wrong")
	ErrPasswordChangeFailed = errors.New("password change failed")
	ErrorNoChange           = errors.New("no change")
	ErrInfoChangeFailed     = errors.New("info change failed")
)
