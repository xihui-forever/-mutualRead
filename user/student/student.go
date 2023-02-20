package student

import (
	"encoding/base64"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/login"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func init() {
	login.LoginHandlerMap[login.LoginTypeAdmin] = func(username interface{}, password string) (uint64, error) {
		data, err := GetStudent(utils.ToUint64(username))
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		err = CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		return data.Id, nil
	}
}

func AddStudent(studentId uint64, pwd string, name string, email string) (*types.ModelStudent, error) {
	a := types.ModelStudent{
		StudentId: studentId,
		Password:  Encrypt(pwd),
		Name:      name,
		Email:     email,
	}

	err := db.Create(&a).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrStudentExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &a, nil
}

func AddStudents(students []types.ModelStudent) (int64, error) {
	var count int64 = 0
	var error error = nil
	var a types.ModelStudent
	for _, value := range students {
		_, err := GetStudent(value.StudentId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				error = ErrStudentExist
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

func RemoveStudent(studentId uint64) (int64, error) {
	var a types.ModelStudent
	result := db.Where("studentId = ?", studentId).Delete(&a)
	err := result.Error
	count := result.RowsAffected
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, ErrStudentNotExist
		}
		log.Errorf("err:%v", err)
		return 0, err
	}
	if count == 0 {
		return 0, ErrStudentRemoveFailed
	}
	return count, nil
}

func RemoveStudents(students []uint64) (int64, error) {
	var count int64 = 0
	for _, value := range students {
		c, err := RemoveStudent(value)
		count += c
		if err != nil {
			log.Errorf("err:%v", err)
			return count, err
		}
	}
	return count, nil
}

func GetStudentsAll() (*[]types.ModelStudent, error) {
	var a []types.ModelStudent
	result := db.Find(&a)
	var err = result.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func GetStudent(studentId uint64) (*types.ModelStudent, error) {
	var a types.ModelStudent
	err := db.Where("studentId = ?", studentId).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStudentNotExist
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

func ChangePassword(studentId uint64, oldPwd, newPwd string) error {
	a, err := GetStudent(studentId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a == nil {
		return ErrStudentNotExist
	}

	if a.Password != Encrypt(oldPwd) {
		return ErrPasswordWrong
	}

	res := db.Model(&types.ModelStudent{}).Where("id = ?", a.Id).Update("password", Encrypt(newPwd))
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

func ChangeEmail(studentId uint64, email string) error {
	a, err := GetStudent(studentId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a.Email == email {
		return ErrorEmailNoChange
	}

	if a.Email != email {
		res := db.Model(&types.ModelStudent{}).Where("id = ?", a.Id).Update("email", email)
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
