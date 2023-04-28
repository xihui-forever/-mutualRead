package main

import (
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/mq"
	"github.com/xihui-forever/mutualRead/appeal"
	_ "github.com/xihui-forever/mutualRead/impl"
	"github.com/xihui-forever/mutualRead/rpc"
	_ "github.com/xihui-forever/mutualRead/teacher"
	"path/filepath"

	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/goon/middleware/storage/redis"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/types"
)

func main() {
	var err error
	log.SetPrefixMsg("mutualRead")

	config.Load()

	err = mq.Start(&mq.Option{
		DataPath:     filepath.Join(utils.GetExecPath(), "data"),
		MemQueueSize: 1,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		types.ModelPerm{},
		types.ModelAppeal{},
		types.ModelPaper{},
		types.ModelExam{},
		types.ModelStudent{},
		types.ModelTeacher{},
		types.ModelAdmin{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	session.SetStorage(redis.New(&redis.RedisConfig{
		Addr: "127.0.0.1:6379",
		DB:   3,
	}))

	err = role.Load()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	goon.WithPermHeader(types.HeaderRoleType)
	rpc.Load()

	err = appeal.Load()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = goon.ListenAndServe(viper.GetString(config.ListenAddr))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
