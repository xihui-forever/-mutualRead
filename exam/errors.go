package exam

import (
	"github.com/xihui-forever/mutualRead/types"
)

var (
	ErrExamNotExist     = types.CreateError(types.ErrExamNotExist)
	ErrExamChangeFailed = types.CreateError(types.ErrExamChangeFailed)
)
