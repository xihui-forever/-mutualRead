package main

import (
	_ "github.com/xihui-forever/mutualRead/impl"
	"github.com/xihui-forever/mutualRead/rpc"
	_ "github.com/xihui-forever/mutualRead/teacher"

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
	log.SetPrefixMsg("mutualRead")

	config.Load()

	err := db.Connect(db.Config{
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

	err = goon.ListenAndServe(viper.GetString(config.ListenAddr))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
