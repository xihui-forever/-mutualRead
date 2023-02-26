package paper

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/image"
	"gorm.io/gorm"
)

func AddPaper(imag image.Derived, grade uint32, examId uint64, examiner uint64, reviewer uint64, chargePerson uint64) (*types.ModelPaper, error) {
	a := types.ModelPaper{
		Imag:     imag,
		Grade:    grade,
		ExamId:   examId,
		Examiner: examiner,
		Reviewer: reviewer,
	}

	err := db.Create(&a).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrPaperExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &a, nil
}

func SetPaper(id uint64, grade uint32) error {
	res := db.Model(&types.ModelPaper{}).Where("id = ?", id).Update("grade", grade)
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

func GetPaper(id uint64) (*types.ModelPaper, error) {
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

func GetPaperListExam(examId uint64) (*[]types.ModelPaper, error) {
	var a []types.ModelPaper
	err := db.Where("exam_id = ?", examId).Find(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}
