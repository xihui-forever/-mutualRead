package teacher

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
		data, err := GetTeacher(utils.ToUint64(username))
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

func AddTeacher(teacherId uint64, pwd string, name string, email string) (*types.ModelTeacher, error) {
	a := types.ModelTeacher{
		TeacherId: teacherId,
		Password:  Encrypt(pwd),
		Name:      name,
		Email:     email,
	}

	err := db.Create(&a).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrTeacherExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &a, nil
}

func GetTeacher(teacherId uint64) (*types.ModelTeacher, error) {
	var a types.ModelTeacher
	err := db.Where("teacherId = ?", teacherId).First(&a).Error
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

func ChangePassword(teacherId uint64, oldPwd, newPwd string) error {
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

func ChangeInfo(teacherId uint64, name string, email string) error {
	a, err := GetTeacher(teacherId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if name == "" || email == "" {
		return ErrTeacherNameOrEmailEmpty
	}

	if a.Name == name && a.Email == email {
		return ErrorNoChange
	}

	if a.Email != email {
		res := db.Model(&types.ModelTeacher{}).Where("id = ?", a.Id).Update("email", email)
		err = res.Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if res.RowsAffected == 0 {
			return ErrInfoChangeFailed
		}
	}
	if a.Name != name {
		res := db.Model(&types.ModelTeacher{}).Where("id = ?", a.Id).Update("name", name)
		err = res.Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if res.RowsAffected == 0 {
			return ErrInfoChangeFailed
		}
	}

	return nil
}

func Encrypt(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(utils.HmacSha384("10e5bbcdadc047328f4ed085ccbf9088", pwd)))
}
