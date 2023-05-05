package admin

import (
	"github.com/xihui-forever/mutualRead/types"
)

var (
	ErrAdminExist           = types.CreateError(types.ErrAdminExist)
	ErrAdminNotExist        = types.CreateError(types.ErrAdminNotExist)
	ErrAdminRemoveFailed    = types.CreateError(types.ErrAdminRemoveFailed)
	ErrPasswordWrong        = types.CreateError(types.ErrPasswordWrong)
	ErrPasswordChangeFailed = types.CreateError(types.ErrPasswordChangeFailed)
)
