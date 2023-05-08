package paper

import (
	"github.com/xihui-forever/mutualRead/types"
)

var (
	ErrPaperExist        = types.CreateError(types.ErrPaperExist)
	ErrPaperChangeFailed = types.CreateError(types.ErrPaperChangeFailed)
	ErrPaperNotExist     = types.CreateError(types.ErrPaperNotExist)
	ErrGradeChangeFailed = types.CreateError(types.ErrGradeChangeFailed)
)
