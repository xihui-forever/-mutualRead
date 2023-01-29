package paper

import "errors"

var (
	ErrPaperExist    = errors.New("paper already exists")
	ErrPaperNotExist = errors.New("paper does not exist")
)
