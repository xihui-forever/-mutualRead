package appeal

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func AddAppeal(appeal *types.ModelAppeal) (*types.ModelAppeal, error) {
	var id uint64
	err := db.Model(&types.ModelAppeal{}).
		Where("paper_id = ?", appeal.PaperId).
		Where(" examiner_id = ?", appeal.ExaminerId).
		Where("state not in (?)", []int{
			types.AppealStateFinish,
			types.AppealStateRecall,
		}).Scan(&id).Error
	if err == gorm.ErrRecordNotFound {

	} else if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	} else {
		return nil, ErrAppealExist
	}

	err = db.Create(&appeal).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return appeal, nil
}

//func RemoveAppealExaminer(id uint64, examinerId uint64) (int64, error) {
//	var a types.ModelAppeal
//	err := db.Where("id = ? AND examiner_id = ?", id, examinerId).First(&a).Error
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return 0, ErrAppealNotExist
//		}
//		log.Errorf("err:%v", err)
//		return 0, err
//	}
//	if a.State != StateWaitReviewer {
//		return 0, ErrRemoveFailed
//	}
//	result := db.Where("id = ? AND state = ?", id, StateWaitReviewer).Delete(&a)
//	err = result.Error
//	count := result.RowsAffected
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return 0, err
//	}
//	if count == 0 {
//		return 0, ErrAppealRemoveFailed
//	}
//	return count, nil
//}
//
//func GetAppealState(id uint64) (int, error) {
//	var a types.ModelAppeal
//	err := db.Where("id = ?", id).First(&a).Error
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return 0, ErrAppealNotExist
//		}
//		log.Errorf("err:%v", err)
//		return 0, err
//	}
//	return a.State, nil
//}
//
//func ChangeToStateWaitTeacher(id uint64) error {
//	res := db.Model(&types.ModelAppeal{}).Where("id = ? AND state = ?", id, StateWaitReviewer).Update("state", StateWaitTeacher)
//	err := res.Error
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//
//	if res.RowsAffected == 0 {
//		return ErrAppealStateChangeFailed
//	}
//
//	return nil
//}
//
//func ChangeToStateProcessed(id uint64) error {
//	res := db.Model(&types.ModelAppeal{}).Where("id = ? and state = ?", id, StateWaitTeacher).Update("state", StateProcessed)
//	err := res.Error
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//
//	if res.RowsAffected == 0 {
//		return ErrAppealStateChangeFailed
//	}
//
//	return nil
//}
//
//func ChangeToStateExpired(id uint64) error {
//	res := db.Model(&types.ModelAppeal{}).Where("id = ? and state = ?", id, StateWaitTeacher).Update("state", StateExpired)
//	err := res.Error
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//
//	if res.RowsAffected == 0 {
//		return ErrAppealStateChangeFailed
//	}
//
//	return nil
//}
//
//func GetAppealsExaminer(examinerId uint64) (*[]types.ModelAppeal, error) {
//	var a []types.ModelAppeal
//	err := db.Where("examiner_id = ?", examinerId).Find(&a).Error
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return nil, nil
//		}
//		log.Errorf("err:%v", err)
//		return nil, err
//	}
//	return &a, nil
//}
//
//func GetAppealsReviewer(reviewerId uint64) (*[]types.ModelAppeal, error) {
//	var a []types.ModelAppeal
//	err := db.Where("reviewer_id = ?", reviewerId).Find(&a).Error
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return nil, nil
//		}
//		log.Errorf("err:%v", err)
//		return nil, err
//	}
//	return &a, nil
//}
//
//func ChangeAppealInfo(id uint64, appealInfo string) error {
//	state, err := GetAppealState(id)
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//	if state == StateWaitReviewer {
//		return ErrInfoCannotChange
//	}
//	res := db.Model(&types.ModelAppeal{}).Where("id = ? AND state = ?", id, StateWaitReviewer).Update("appeal_info", appealInfo)
//	err = res.Error
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//
//	if res.RowsAffected == 0 {
//		return ErrAppealInfoChangeFailed
//	}
//
//	return nil
//}
//
//func ChangeReviewInfo(id uint64, reviewInfo string) error {
//	state, err := GetAppealState(id)
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//	if state != StateWaitReviewer {
//		return ErrReviewInfoCannotChange
//	}
//	res := db.Model(&types.ModelAppeal{}).Where("id = ? AND state = ?", id, StateWaitReviewer).Update("review_info", reviewInfo)
//	err = res.Error
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//
//	if res.RowsAffected == 0 {
//		return ErrReviewInfoChangeFailed
//	}
//
//	return nil
//}
//
//func ChangeAppealResult(id uint64, appealResult string) error {
//	state, err := GetAppealState(id)
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//	if state != StateWaitTeacher {
//		return ErResultCannotChange
//	}
//	res := db.Model(&types.ModelAppeal{}).Where("id = ? AND state = ?", id, StateWaitTeacher).Update("appealResult", appealResult)
//	err = res.Error
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return err
//	}
//
//	if res.RowsAffected == 0 {
//		return ErrAppealResultChangeFailed
//	}
//
//	return nil
//}

func ListAppeal(opts *types.ListOption) ([]*types.ModelAppeal, *types.Page, error) {
	db := db.Model(&types.ModelAppeal{})

	for _, option := range opts.Options {
		switch option.Key {
		case types.ListAppeal_OptionTeacherId:
			db = db.Where("teacher_id = ?", option.Val)
		case types.ListAppeal_OptionReviewerId:
			db = db.Where("reviewer_id = ?", option.Val)
		case types.ListAppeal_OptionExaminerId:
			db = db.Where("examiner_id = ?", option.Val)
		case types.ListAppeal_OptionPaperId:
			db = db.Where("paper_id = ?", option.Val)
		case types.ListAppeal_OptionId:
			db = db.Where("id = ?", option.Val)
		case types.ListAppeal_OptionStates:
			db = db.Where("state in (?)", option.Val)
		}
	}

	var list []*types.ModelAppeal
	page, err := opts.Find(db, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	return list, page, nil
}
