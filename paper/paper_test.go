package paper

import (
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/types"
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

	a := types.ModelPaper{
		Img:        "C:\\Users\\Lenovo\\Pictures\\idle.png",
		Grade:      60,
		ExamId:     1,
		ExaminerId: 20190001,
		ReviewerId: 20190002,
	}

	var paper *types.ModelPaper
	paper, err = Add(a)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(paper)
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

	res, err := GetPaperListExaminer(20190001)
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

	res, err := GetPaperListExam(1, 9112019111)
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

	res, err := GetPaper(1, 9112019111)
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

	err = ChangePaperGrade(1, 99, 9112019111)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
