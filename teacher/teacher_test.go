package teacher

import (
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/goon/config"
	"github.com/xihui-forever/goon/types"
	"testing"
)

func TestAddTeacher(t *testing.T) {
	config.Load()
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelTeacher{}, // db.AutoMigrate(&types.ModelAdmin{})
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	a, err := AddTeacher("9112019111", "123456", "老师", "1234567890@123.com")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(a)
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

	err = ChangePassword("9112019111", "123456", "654321")
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

	err = ChangeInfo("9112019111", "laoshi", "123@123.com")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
