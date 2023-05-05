package student

import (
	"github.com/xihui-forever/mutualRead/types"
)

var (
	ErrStudentExist         = types.CreateError(types.ErrStudentExist)
	ErrStudentNotExist      = types.CreateError(types.ErrStudentNotExist)
	ErrStudentRemoveFailed  = types.CreateError(types.ErrStudentRemoveFailed)
	ErrPasswordWrong        = types.CreateError(types.ErrPasswordWrong)
	ErrPasswordChangeFailed = types.CreateError(types.ErrPasswordChangeFailed)
	ErrorEmailNoChange      = types.CreateError(types.ErrorEmailNoChange)
	ErrEmailChangeFailed    = types.CreateError(types.ErrEmailChangeFailed)
)
