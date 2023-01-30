package appeal

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func AddAppeal(paperId uint64) (*types.ModelAppeal, error) {
	a := types.ModelAppeal{
		PaperId:      paperId,
		AppealStatus: "等待阅卷人回应",
	}

	err := db.Create(&a).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrAppealExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &a, nil
}

func GetAppealsByPaperId(paperId uint64) (*[]types.ModelAppeal, error) {
	var a []types.ModelAppeal
	err := db.Where("paperId = ?", paperId).Find(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//这儿应该不返回一个错误吧
			return nil, nil
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func ChangeAppealInfo(id uint64, appealInfo string) error {
	res := db.Model(&types.ModelAppeal{}).Where("id = ?", id).Update("appealInfo", appealInfo)
	err := res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if res.RowsAffected == 0 {
		return ErrAppealInfoChangeFailed
	}

	return nil
}

func ChangeReviewInfo(id uint64, reviewInfo string) error {
	res := db.Model(&types.ModelAppeal{}).Where("id = ?", id).Update("reviewInfo", reviewInfo)
	err := res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if res.RowsAffected == 0 {
		return ErrReviewInfoChangeFailed
	}

	return nil
}

func ChangeAppealResult(id uint64, appealResult string) error {
	res := db.Model(&types.ModelAppeal{}).Where("id = ?", id).Update("appealResult", appealResult)
	err := res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if res.RowsAffected == 0 {
		return ErrAppealResultChangeFailed
	}

	return nil
}
