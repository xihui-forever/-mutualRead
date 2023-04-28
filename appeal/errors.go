package appeal

import "errors"

var (
	ErrAppealFailed             = errors.New("one appeal already progress")
	ErrAppealExist              = errors.New("appeal already exists")
	ErrAppealNotExist           = errors.New("appeal not exists")
	ErrRemoveFailed             = errors.New("appeal cannot be removed")
	ErrAppealRemoveFailed       = errors.New("appeal remove failed")
	ErrInfoCannotChange         = errors.New("info cannot be changed")
	ErrAppealStateChangeFailed  = errors.New("appeall state change failed")
	ErrAppealInfoChangeFailed   = errors.New("appeal info change failed")
	ErrReviewInfoChangeFailed   = errors.New("review info change failed")
	ErrReviewInfoCannotChange   = errors.New("review info cannot be changed")
	ErResultCannotChange        = errors.New("result cannot be changed")
	ErrAppealResultChangeFailed = errors.New("appeal result change failed")

	ErrAppealAlreadyHanded = errors.New("申诉已处理")
)
