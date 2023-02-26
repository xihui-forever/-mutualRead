package appeal

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/types"
	"testing"
)

func TestAddAppeal(t *testing.T) {
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

	var p *types.ModelPaper
	p, err = paper.GetPaper(1)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	a := types.ModelAppeal{
		State:      StateWaitReviewer,
		PaperId:    1,
		ExaminerId: p.Examiner,
		ReviewerId: p.Reviewer,
		AppealInfo: "第一题误判",
	}
	var appeal *types.ModelAppeal
	appeal, err = AddAppeal(a)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(appeal)
}

func TestChangeAppealInfo(t *testing.T) {
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

	err = ChangeAppealInfo(1, "第一题第一小问误判")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestChangeReviewInfo(t *testing.T) {
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

	err = ChangeReviewInfo(1, "误判")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}

func TestChangeAppealResult(t *testing.T) {
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

	err = ChangeAppealResult(1, "考试人成绩加五分，阅卷人成绩减五分")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
