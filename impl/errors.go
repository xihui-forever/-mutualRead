package impl

import "errors"

var (
	ErrPasswordIncorrect = errors.New("password is incorrect")
	ErrTeacherNotExist   = errors.New("teacher not exist")
)
