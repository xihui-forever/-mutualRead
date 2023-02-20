package exam

import "errors"

var (
	ErrExamNotExist     = errors.New("exam not found")
	ErrExamChangeFailed = errors.New("exam name change failed")
)
