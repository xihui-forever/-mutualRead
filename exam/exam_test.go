package exam

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/types"
	"testing"
)

func TestAddExam(t *testing.T) {
	config.Load()
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelAppeal{},
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	var a *types.ModelExam
	a, err = Add("数据结构第一次测试", "90120001")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	t.Log(a)
}
