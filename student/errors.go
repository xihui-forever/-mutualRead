package student

import "errors"

var (
	ErrStudentExist         = errors.New("student already exists")
	ErrStudentNotExist      = errors.New("student not exists")
	ErrStudentRemoveFailed  = errors.New("student remove failed")
	ErrPasswordWrong        = errors.New("password is wrong")
	ErrPasswordChangeFailed = errors.New("password change failed")
	ErrorEmailNoChange      = errors.New("email no change")
	ErrEmailChangeFailed    = errors.New("email change failed")
)
