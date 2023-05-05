package appeal

import (
	"github.com/xihui-forever/mutualRead/types"
)

var (
	ErrAppealExist    = types.CreateError(types.ErrAppealAlreadyExists)
	ErrAppealNotExist = types.CreateError(types.ErrAppealNotExists)

	ErrAppealAlreadyHanded = types.CreateError(types.ErrAppealAlreadyHanded)
)
