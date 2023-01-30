package paper

import "errors"

var (
	ErrPaperExist        = errors.New("paper already exists")
	ErrPaperNotExist     = errors.New("paper does not exist")
	ErrGradeChangeFailed = errors.New("grade change failed")
)
