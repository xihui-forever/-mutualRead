package paper

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func Add(paper *types.ModelPaper) (*types.ModelPaper, error) {
	err := db.Create(&paper).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrPaperExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return paper, nil
}

func ChangePaperGrade(id uint64, grade uint32, teacherId uint64) error {
	res := db.Model(&types.ModelPaper{}).Where("id = ? AND teacher_id = ?", id, teacherId).Update("grade", grade)
	err := res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if res.RowsAffected == 0 {
		return ErrPaperChangeFailed
	}
	return nil
}

func GetPaper(id uint64, teacherId uint64) (*types.ModelPaper, error) {
	var a types.ModelPaper
	err := db.Where("id = ? AND teacher_id = ?", id, teacherId).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func Get(id uint64) (*types.ModelPaper, error) {
	var a types.ModelPaper
	err := db.Where("id = ?", id).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetPaperListAdmin() (*[]types.ModelPaper, error) {
	var a []types.ModelPaper
	result := db.Find(&a)
	var err = result.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetPaperListExaminer(examiner uint64) (*[]types.ModelPaper, error) {
	var a []types.ModelPaper
	err := db.Where("examiner = ?", examiner).Find(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetPaperListExam(examId uint64, teacherId uint64) (*[]types.ModelPaper, error) {
	var a []types.ModelPaper
	err := db.Where("exam_id = ? AND teacher_id = ?", examId, teacherId).Find(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func ListPaper(opts *types.ListOption) ([]*types.ModelPaper, *types.Page, error) {
	db := db.Model(&types.ModelPaper{})

	for _, option := range opts.Options {
		switch option.Key {
		case types.ListPaper_OptionExamId:
			db = db.Where("exam_id = ?", option.Val)
		case types.ListPaper_OptionExaminerId:
			db = db.Where("examiner = ?", option.Val)
		case types.ListPaper_OptionReviewerId:
			db = db.Where("reviewer = ?", option.Val)
		case types.ListPaper_OptionTeacherId:
			db = db.Where("teacher_id = ?", option.Val)
		}
	}

	var list []*types.ModelPaper
	page, err := opts.Find(db, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	return list, page, nil
}

func DelPaper(id uint64, teacherId uint64) error {
	err := db.Where("id = ?", id).Where("teacher_id = ?", teacherId).Delete(&types.ModelPaper{}).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
