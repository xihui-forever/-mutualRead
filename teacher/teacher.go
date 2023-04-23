package teacher

import (
	"encoding/base64"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func AddTeacher(teacher types.ModelTeacher) (*types.ModelTeacher, error) {
	err := db.Create(&teacher).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrTeacherExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &teacher, nil
}

func AddTeachers(teachers []types.ModelTeacher) (int64, error) {
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

func GetTeacherById(id uint64) (*types.ModelTeacher, error) {
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

func ChangeEmail(teacherId string, email string) error {
	a, err := GetTeacher(teacherId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if email == "" {
		return ErrEmailEmpty
	}

	if a.Email == email {
		return ErrorEmailNoChange
	}

	if a.Email != email {
		res := db.Model(&types.ModelTeacher{}).Where("id = ?", a.Id).Update("email", email)
		err = res.Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if res.RowsAffected == 0 {
			return ErrEmailChangeFailed
		}
	}

	return nil
}

func Encrypt(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(utils.HmacSha384("10e5bbcdadc047328f4ed085ccbf9088", pwd)))
}
