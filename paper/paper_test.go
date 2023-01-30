package paper

import (
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/types"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/image"
	"testing"
)

func TestAddPaper(t *testing.T) {
	config.Load()
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

	a, err := AddPaper("软件工程01", image.Derived{
		Fingerprint:     nil,
		Distance:        0,
		LayerInfo:       nil,
		BaseResourceUrl: "",
	}, 95, 20199999, 20190000, 9112019111)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(a)
}

func TestGetPapersByExaminer(t *testing.T) {
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

	res, err := GetPapersByExaminer(20199999)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(res)
}

func TestGetPapersByName(t *testing.T) {
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

	res, err := GetPapersByName("软件工程01")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(res)
}

func TestGetPaperById(t *testing.T) {
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

	res, err := GetPaperById(1)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(res)
}

func TestChangeGrade(t *testing.T) {
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

	err = ChangeGrade(1, 99)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
