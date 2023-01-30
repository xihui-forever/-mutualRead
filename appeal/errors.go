package appeal

import "errors"

var (
	ErrAppealExist              = errors.New("appeal already exists")
	ErrAppealNotExist           = errors.New("appeal not exists")
	ErrAppealInfoChangeFailed   = errors.New("appeal info change failed")
	ErrReviewInfoChangeFailed   = errors.New("review info change failed")
	ErrAppealResultChangeFailed = errors.New("appeal result change failed")
)
