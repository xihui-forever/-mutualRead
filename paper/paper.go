package paper

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/types"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/image"
	"gorm.io/gorm"
)

func init() {
	role.AddRole("admin", AddPaper)
}

func AddPaper(name string, imag image.Derived, grade uint32, examiner uint64, reviewer uint64, chargePerson uint64) (*types.ModelPaper, error) {
	a := types.ModelPaper{
		Name:         name,
		Imag:         imag,
		Grade:        grade,
		Examiner:     examiner,
		Reviewer:     reviewer,
		ChargePerson: chargePerson,
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

func GetPapersByExaminer(examiner uint64) (*[]types.ModelPaper, error) {
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

func GetPapersByName(name string) (*[]types.ModelPaper, error) {
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

func GetPaperById(id uint64) (*types.ModelPaper, error) {
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

func ChangeGrade(id uint64, grade uint32) error {
	res := db.Model(&types.ModelPaper{}).Where("id = ?", id).Update("grade", grade)
	err := res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if res.RowsAffected == 0 {
		return ErrGradeChangeFailed
	}

	return nil
}
