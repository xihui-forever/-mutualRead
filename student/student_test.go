package student

import (
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/types"
	"testing"
)

func TestAddStudent(t *testing.T) {
	config.Load()
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelStudent{},
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	a := types.ModelStudent{
		StudentId: "20190001",
		Password:  Encrypt("111"),
		Name:      "张三",
		Email:     "email@123.com",
	}
	stu, err := Add(a)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(stu)
}
