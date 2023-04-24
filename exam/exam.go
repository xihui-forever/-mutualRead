package exam

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
)

func AddExam(name string, teacherId string) (*types.ModelExam, error) {
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
		log.Errorf("err:%v", err)
		return err
	}
	return nil
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

func GetExam(id uint64) (*types.ModelExam, error) {
	var a types.ModelExam
	err := db.Where("id = ?", id).First(&a).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetExamList(opt *types.ListOption) ([]*types.ModelExam, error) {
	var a []*types.ModelExam

	//db := types.GetDbFromListOption(opt)

	db := db.Model(&types.ModelExam{})

	for _, option := range opt.Options {
		switch option.Key {
		case types.ExamListReq_OptionTeacherId:
			db.Where("teacher_id = ?", option.Val)
		default:
			return nil, types.CreateError(types.ErrInvalidOptionKey)
		}
	}

	err := db.Find(&a).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return a, nil
}
