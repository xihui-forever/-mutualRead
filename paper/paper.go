package paper

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/image"
	"gorm.io/gorm"
)

func AddPaper(name string, imag image.Derived, grade uint32, examiner uint64, reviewer uint64) (*types.ModelPaper, error) {
	a := types.ModelPaper{
		Name:     name,
		Imag:     imag,
		Grade:    grade,
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

func GetPaperByExaminer(examiner uint32) (*types.ModelPaper, error) {
	var a types.ModelPaper
	err := db.Where("examiner = ?", examiner).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetPaperByReviewer(reviewer uint32) (*types.ModelPaper, error) {
	var a types.ModelPaper
	err := db.Where("reviewer = ?", reviewer).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetPapers(name string) (*[]types.ModelPaper, error) {
	var a []types.ModelPaper
	err := db.Where("name = ?", name).Find(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrPaperNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}
