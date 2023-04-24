package student

import (
	"encoding/base64"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
)

func AddStudent(student types.ModelStudent) (*types.ModelStudent, error) {
	err := db.Create(&student).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrStudentExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &student, nil
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

func RemoveStudent(studentId string) (int64, error) {
	var a types.ModelStudent
	result := db.Where("studentId = ?", studentId).Delete(&a)
	err := result.Error
	count := result.RowsAffected
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, err
	}
	if count == 0 {
		return 0, ErrStudentRemoveFailed
	}
	return count, nil
}

func RemoveStudents(students []string) (int64, error) {
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

func GetStudent(studentId string) (*types.ModelStudent, error) {
	var a types.ModelStudent
	err := db.Where("student_id = ?", studentId).First(&a).Error
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

func ChangePassword(studentId string, oldPwd, newPwd string) error {
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

func ChangeEmail(id uint64, email string) error {
	err := db.Where("id = ?", id).Update("email", email).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func Encrypt(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(utils.HmacSha384("10e5bbcdadc047328f4ed085ccbf9088", pwd)))
}

func Get(id uint64) (*types.ModelStudent, error) {
	var a types.ModelStudent
	err := db.Where("id = ?", id).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStudentNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &a, nil
}

func Set(s *types.ModelStudent) (*types.ModelStudent, error) {
	err := db.Model(&types.ModelStudent{}).Omit("password").Save(s).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return s, nil
}

func Del(id uint64) error {
	err := db.Where("id = ?", id).Delete(&types.ModelStudent{}).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func List(opt *types.ListOption) ([]*types.ModelStudent, *types.Page, error) {
	db := db.Model(&types.ModelStudent{})

	for _, option := range opt.Options {
		switch option.Key {
		case types.ListStudent_OptionStudentId:
			db = db.Where("student_id = ?", option.Val)
		case types.ListStudent_OptionNameLike:
			db = db.Where("name like ?", "%"+option.Val+"%")
		case types.ListStudent_OptionEmailLike:
			db = db.Where("email like ?", "%"+option.Val+"%")
		default:
			log.Errorf("unknown option:%v", option)
		}
	}

	var list []*types.ModelStudent
	page, err := opt.Find(db, &list)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, nil, err
	}

	return list, page, nil
}
