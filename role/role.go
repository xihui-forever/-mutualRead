package role

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func Load() error {
	_, err := BatchAddRolePerm(types.RoleTypeAdmin, []string{"/exam_add"})
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}
	_, err = BatchAddRolePerm(types.RoleTypeTeacher, []string{
		types.CmdPathGetTeacher,
		types.CmdPathChangeTeacher,
		"/email_update", "/password_update",
		"/exam_add", "/exam_update", "/exam_list", "/exam_delete", "/exam_get",
		"/paper_add", "/paper_delete", "/paper_list", "/paper_get"})
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}
	_, err = BatchAddRolePerm(types.RoleTypePublic, []string{"/login"})
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

func BatchAddRolePerm(role int, permissions []string) (int64, error) {
	var count int64 = 0
	var error error = nil
	for _, value := range permissions {
		a := types.ModelPerm{
			Role:       role,
			Permission: value,
		}

		exist, err := CheckPermission(role, value)
		if err != nil {
			log.Errorf("err:%s", err)
			return 0, err
		}

		if exist {
			continue
		}

		err = db.Create(&a).Error
		if err != nil {
			if !types.IsUniqueErr(err) {
				log.Errorf("err:%v", err)
				return count, err
			}
		}
		count++
	}
	return count, error
}

func CheckPermission(role int, permission string) (bool, error) {
	var id uint64
	err := db.Model(&types.ModelPerm{}).Where("role = ? AND permission = ?", role, permission).Select("id").Scan(&id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, ErrRolePermExists
		}
		log.Errorf("err:%v", err)
		return false, err
	}

	return id > 0, nil
}
