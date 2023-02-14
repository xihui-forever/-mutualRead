package teacher

import "errors"

var (
	ErrTeacherExist            = errors.New("Teacher already exists")
	ErrTeacherNotExist         = errors.New("Teacher does not exist")
	ErrorNewPwdEmpty           = errors.New("New password is null")
	ErrPasswordWrong           = errors.New("Password is wrong")
	ErrPasswordChangeFailed    = errors.New("Password change failed")
	ErrInfoChangeFailed        = errors.New("Info change failed")
	ErrorNoChange              = errors.New("New Info equals existing info")
	ErrTeacherNameOrEmailEmpty = errors.New("Teacher name or email must be not null")
)
