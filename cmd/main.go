package main

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/impl"
	"github.com/xihui-forever/mutualRead/role"
)

func main() {
	config.Load()

	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	app := goon.New()

	app.PreUse("/", func(ctx *goon.Ctx) (string, error) {
		var sess *session.Session
		sid := ctx.GetReqHeader("X-Session-Id")
		s, err := sess.GetSession(sid)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}
		return s, nil
	})

	app.PreUse("/", func(ctx *goon.Ctx) error {
		roleType := utils.ToInt(ctx.GetReqHeader("role"))
		path := ctx.Path()
		_, err := role.CheckPermission(roleType, path)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})

	for _, value := range impl.CmdList {
		app.Post(value.Path, value.Logic)
	}

	_ = fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		err := app.Call(ctx)
		if err != nil {
			log.Errorf("err:%v", err)
			ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
			_, e := ctx.Write([]byte(err.Error()))
			if e != nil {
				log.Errorf("err:%v", e)
			}
		}
	})
}
