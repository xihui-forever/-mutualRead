package teacher

import (
	"github.com/xihui-forever/mutualRead/types"
)

var (
	ErrTeacherExist         = types.CreateError(types.ErrTeacherExist)
	ErrTeacherNotExist      = types.CreateError(types.ErrTeacherNotExist)
	ErrTeacherRemoveFailed  = types.CreateError(types.ErrTeacherRemoveFailed)
	ErrorNewPwdEmpty        = types.CreateError(types.ErrorNewPwdEmpty)
	ErrPasswordWrong        = types.CreateError(types.ErrPasswordWrong)
	ErrPasswordChangeFailed = types.CreateError(types.ErrPasswordChangeFailed)
	ErrEmailChangeFailed    = types.CreateError(types.ErrEmailChangeFailed)
	ErrorEmailNoChange      = types.CreateError(types.ErrorEmailNoChange)
	ErrEmailEmpty           = types.CreateError(types.ErrEmailEmpty)
)
