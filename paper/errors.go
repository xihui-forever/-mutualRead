package paper

import "errors"

var (
	ErrPaperExist        = errors.New("paper already exists")
	ErrPaperChangeFailed = errors.New("paper change failed")
	ErrPaperNotExist     = errors.New("paper does not exist")
	ErrGradeChangeFailed = errors.New("grade change failed")
)
