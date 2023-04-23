package teacher

import (
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/types"
	"testing"
)

func TestAddTeacher(t *testing.T) {
	config.Load()
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelTeacher{},
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	a := types.ModelTeacher{
		TeacherId: "90120001",
		Password:  Encrypt("123"),
		Name:      "李强",
		Email:     "1234567890@123.com",
	}
	teacher, err := AddTeacher(a)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(teacher)
}

func TestChangePassword(t *testing.T) {
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelTeacher{},
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	err = ChangePassword("90120001", "123456", "654321")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestChangeInfo(t *testing.T) {
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelTeacher{},
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	err = ChangeEmail("90120001", "123@123.com")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestGetTeacher(t *testing.T) {
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelPaper{},
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	res, err := GetTeacher("90120001")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(res)
}
