package impl

import (
	"github.com/xihui-forever/mutualRead/types"
)

var (
	ErrPasswordIncorrect = types.CreateError(types.ErrPasswordIncorrect)
	ErrTeacherNotExist   = types.CreateError(types.ErrTeacherNotExist)
)
