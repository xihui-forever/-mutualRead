package exam

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func Add(a *types.ModelExam) (*types.ModelExam, error) {
	err := db.Create(&a).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return a, nil
}

func ChangeExamName(id uint64, name string) error {
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

func Get(id uint64) (*types.ModelExam, error) {
	var a types.ModelExam
	err := db.Where("id = ?", id).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrExamNotExist
		}

		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func Set(e *types.ModelExam) (*types.ModelExam, error) {
	err := db.Save(&e).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return e, nil
}

func Del(id uint64) error {
	err := db.Where("id = ?", id).Delete(&types.ModelExam{}).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func List(opt *types.ListOption) ([]*types.ModelExam, *types.Page, error) {
	db := db.Model(&types.ModelExam{})

	for _, option := range opt.Options {
		switch option.Key {
		case types.ListExam_OptionTeacherId:
			db = db.Where("teacher_id = ?", option.Val)
		case types.ListExam_OptionNameLike:
			db = db.Where("name like ?", "%"+option.Val+"%")
		}
	}

	var exams []*types.ModelExam
	page, err := opt.Find(db, &exams)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	return exams, page, nil
}
