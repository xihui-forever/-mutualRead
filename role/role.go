package role

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func init() {
	BatchAddRolePerm(1, []string{"/exam_add"})
}

func BatchAddRolePerm(role int, permissions []string) (int64, error) {
	var count int64 = 0
	var error error = nil
	for _, value := range permissions {
		a := types.ModelPerm{
			Role:       role,
			Permission: value,
		}

		result := db.Create(&a)
		if result.Error != nil {
			if types.IsUniqueErr(result.Error) {
				error = ErrRolePermExists
			} else {
				log.Errorf("err:%v", result.Error)
				return count, result.Error
			}
		}
		count += result.RowsAffected
	}
	return count, error
}

func CheckPermission(role int, permission string) (bool, error) {
	var a types.ModelPerm
	err := db.Where("role = ? AND permission = ?", role, permission).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, ErrPermNotExist
		}
		log.Errorf("err:%v", err)
		return false, err
	}
	return true, nil
}
