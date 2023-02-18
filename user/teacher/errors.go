package teacher

import "errors"

var (
	ErrTeacherExist         = errors.New("teacher already exists")
	ErrTeacherNotExist      = errors.New("teacher does not exist")
	ErrTeacherRemoveFailed  = errors.New("teacher remove failed")
	ErrorNewPwdEmpty        = errors.New("new password is null")
	ErrPasswordWrong        = errors.New("password is wrong")
	ErrPasswordChangeFailed = errors.New("password change failed")
	ErrEmailChangeFailed    = errors.New("email change failed")
	ErrorEmailNoChange      = errors.New("new Email equals existing info")
	ErrEmailEmpty           = errors.New("teacher email must be not null")
)
