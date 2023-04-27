package teacher

import (
	"encoding/base64"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func Add(teacher *types.ModelTeacher) (*types.ModelTeacher, error) {
	if teacher.Password == "" {
		return nil, types.CreateErrorWithMsg(types.ErrInvalidParam, "未填写密码")
	}

	teacher.Password = Encrypt(teacher.Password)

	err := db.Create(&teacher).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrTeacherExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return teacher, nil
}

func Set(t *types.ModelTeacher) (*types.ModelTeacher, error) {
	err := db.Model(&types.ModelTeacher{}).Omit("password").Save(&t).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return t, nil
}

func AddTeachers(teachers []*types.ModelTeacher) (int64, error) {
	var count int64 = 0
	var error error = nil
	var a types.ModelTeacher
	for _, value := range teachers {
		_, err := GetTeacher(value.TeacherId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				error = ErrTeacherExist
				count += 0
			} else {
				log.Errorf("err:%v", err)
				return count, err
			}
		} else {
			result := db.Create(&a)
			if result.Error != nil {
				log.Errorf("err:%v", result.Error)
				return count, result.Error
			}
			count += result.RowsAffected
		}
	}
	return count, error
}

func RemoveTeacher(teacherId string) (int64, error) {
	var a types.ModelTeacher
	result := db.Where("teacher_id = ?", teacherId).Delete(&a)
	err := result.Error
	count := result.RowsAffected
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, err
	}
	if count == 0 {
		return 0, ErrTeacherRemoveFailed
	}
	return count, nil
}

func RemoveTeachers(teachers []string) (int64, error) {
	var count int64 = 0
	for _, value := range teachers {
		c, err := RemoveTeacher(value)
		count += c
		if err != nil {
			log.Errorf("err:%v", err)
			return count, err
		}
	}
	return count, nil
}

func GetTeachersAll() (*[]types.ModelTeacher, error) {
	var a []types.ModelTeacher
	result := db.Find(&a)
	var err = result.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func Get(id uint64) (*types.ModelTeacher, error) {
	var a types.ModelTeacher
	err := db.Where("id = ?", id).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTeacherNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetTeacher(teacherId string) (*types.ModelTeacher, error) {
	var a types.ModelTeacher
	err := db.Where("teacher_id = ?", teacherId).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTeacherNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func CheckPassword(input string, password string) error {
	if Encrypt(input) != password {
		return ErrPasswordWrong
	}
	return nil
}

func ChangePassword(teacherId string, oldPwd, newPwd string) error {
	if newPwd == "" {
		return ErrorNewPwdEmpty
	}

	a, err := GetTeacher(teacherId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a == nil {
		return ErrTeacherNotExist
	}

	if a.Password != Encrypt(oldPwd) {
		return ErrPasswordWrong
	}

	res := db.Model(&types.ModelTeacher{}).Where("id = ?", a.Id).Update("password", Encrypt(newPwd))
	err = res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if res.RowsAffected == 0 {
		return ErrPasswordChangeFailed
	}

	return nil
}

func ChangeEmail(id uint64, email string) error {
	err := db.Model(&types.ModelTeacher{}).Where("id = ?", id).Update("email", email).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func Encrypt(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(utils.HmacSha384("10e5bbcdadc047328f4ed085ccbf9088", pwd)))
}

func ListTeacher(opt *types.ListOption) ([]*types.ModelTeacher, *types.Page, error) {
	db := db.Model(&types.ModelTeacher{})

	for _, option := range opt.Options {
		switch option.Key {
		case types.ListTeacher_OptionTeacherId:
			db = db.Where("teacher_id = ?", option.Val)
		case types.ListTeacher_OptionNameLike:
			db = db.Where("name like %?%", option.Val)
		case types.ListTeacher_OptionEmailLike:
			db = db.Where("email like %?%", option.Val)
		default:
			log.Errorf("unknown option key:%v", option.Key)
		}
	}

	var teachers []*types.ModelTeacher
	page, err := opt.Find(db, &teachers)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	return teachers, page, nil
}

func DelTeacher(id uint64) error {
	err := db.Where("id = ?", id).Delete(&types.ModelTeacher{}).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func ResetPassword(username, password string) error {
	var a types.ModelTeacher
	err := db.Where("teacher_id = ?", username).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrTeacherNotExist
		}
		log.Errorf("err:%v", err)
		return err
	}

	err = db.Model(&a).Where("id = ?", a.Id).Update("password", Encrypt(password)).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
