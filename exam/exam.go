package exam

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func AddExam(name string, teacherId uint64) (*types.ModelExam, error) {
	a := types.ModelExam{
		Name:      name,
		TeacherId: teacherId,
	}

	err := db.Create(&a).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &a, nil
}

func RemoveExam(id uint64) error {
	var a types.ModelExam
	err := db.Where("id = ?", id).Delete(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrExamNotExist
		}
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func SetExam(id uint64, name string) error {
	var a types.ModelExam
	result := db.Model(&types.ModelExam{}).Where("id = ?").First(&a).Update("name", "name")
	err := result.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if result.RowsAffected == 0 {
		return ErrExamChangeFailed
	}

	return nil
}

func GetExam(id uint64) (*types.ModelExam, error) {
	var a types.ModelExam
	err := db.Where("id = ?", id).First(&a).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetExamList(teacherId uint64) (*[]types.ModelExam, error) {
	var a []types.ModelExam
	err := db.Where("teacherId = ?", teacherId).Find(&a).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}
