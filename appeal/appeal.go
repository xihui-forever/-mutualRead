package appeal

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

const (
	StateWaitReviewer = iota + 1
	StateWaitTeacher
	StateProcessed
	StateExpired
)

func AddAppeal(appeal types.ModelAppeal) (*types.ModelAppeal, error) {
	p, err := paper.GetPaper(appeal.PaperId)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	var a []types.ModelAppeal
	res := db.Where("paper_id = ?", p.Id).Find(&a)
	if res.Error != nil {
		log.Errorf("err:%v", err)
		return nil, res.Error
	}
	for _, value := range a {
		if value.State != StateProcessed || value.State != StateExpired {
			return nil, ErrAppealFailed
		}
	}
	err = db.Create(&appeal).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrAppealExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &appeal, nil
}

func RemoveAppealExaminer(id uint64) (int64, error) {
	var a types.ModelAppeal
	err := db.Where("appeal_id = ?", id).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, ErrAppealNotExist
		}
		log.Errorf("err:%v", err)
		return 0, err
	}
	if a.State != StateWaitReviewer {
		return 0, ErrRemoveFailed
	}
	result := db.Where("appeal_id = ?", id).Delete(&a)
	err = result.Error
	count := result.RowsAffected
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, err
	}
	if count == 0 {
		return 0, ErrAppealRemoveFailed
	}
	return count, nil
}

func GetAppealState(id uint64) (int, error) {
	var a types.ModelAppeal
	err := db.Where("appeal_id = ?", id).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, ErrAppealNotExist
		}
		log.Errorf("err:%v", err)
		return 0, err
	}
	return a.State, nil
}

func GetAppealsExaminer(examinerId uint64) (*[]types.ModelAppeal, error) {
	var a []types.ModelAppeal
	err := db.Where("examiner_id = ?", examinerId).Find(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetAppealsReviewer(reviewerId uint64) (*[]types.ModelAppeal, error) {
	var a []types.ModelAppeal
	err := db.Where("reviewer_id = ?", reviewerId).Find(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func ChangeAppealInfo(id uint64, appealInfo string) error {
	state, err := GetAppealState(id)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if state != StateWaitReviewer {
		return ErrInfoCannotChange
	}
	res := db.Model(&types.ModelAppeal{}).Where("id = ?", id).Update("appeal_info", appealInfo)
	err = res.Error
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
	state, err := GetAppealState(id)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if state != StateWaitReviewer {
		return ErrReviewInfoCannotChange
	}
	res := db.Model(&types.ModelAppeal{}).Where("id = ?", id).Update("review_info", reviewInfo)
	err = res.Error
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
	state, err := GetAppealState(id)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if state != StateWaitTeacher {
		return ErResultCannotChange
	}
	res := db.Model(&types.ModelAppeal{}).Where("id = ?", id).Update("appealResult", appealResult)
	err = res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if res.RowsAffected == 0 {
		return ErrAppealResultChangeFailed
	}

	return nil
}
