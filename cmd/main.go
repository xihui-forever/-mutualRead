package main

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/goon/middleware/storage/redis"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/login"
	"github.com/xihui-forever/mutualRead/types"
	"net/http"
)

func main() {
	config.Load()

	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelAdmin{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	mux := goon.New()
	sess := session.New(redis.New(&redis.RedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	}))

	mux.Post("/admin/login", func(ctx *handler.Ctx, req *login.LoginReq) (*login.LoginRsp, error) {
		resp, err := login.LoginHandler(req, sess)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		ctx.SetSid(resp.Token)
		return resp, nil
	})

	mux.PreUse("/user", func(ctx *handler.Ctx) error {
		data, err := session.GetSession(ctx.GetSid())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		log.Infof("data:%v", data)

		return nil
	})

	_ = http.ListenAndServe(":8080", mux)
}
