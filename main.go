package main

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/appeal"
	"github.com/xihui-forever/mutualRead/cmd"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/types"
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

	for _, value := range cmd.CmdList {
		app.PreUse(value.Path, func(ctx *goon.Ctx) error {
			/*session.Handler(&session.Option[]{
				Session:    session.New(redis.New(&redis.RedisConfig{
					Addr:     "127.0.0.1:6379",
					Password: "",
					DB:       1,
				})),
				Header:     "X-Session-Id",
				Expiration: 0,
				NeedSkip:   nil,
				OnError:    nil,
				OnSuccess:  nil,
			})*/
			var sess *session.Session
			sid := ctx.GetReqHeader("X-Session-Id")
			_, err := sess.GetSession(sid)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			return nil
		})

		app.PreUse(value.Path, func(value cmd.Cmd) error {
			_, err := role.CheckPermission(value.Role, value.Path)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			return nil
		})
	}

	app.Post("/appeal_add", func(ctx *goon.Ctx) (*types.ModelAppeal, error) {
		//获取paper_id,teacherId,info
		var paperId uint64
		var teacherId uint64
		var info string

		var p *types.ModelPaper
		p, err = paper.GetPaper(paperId, teacherId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		a := types.ModelAppeal{
			State:      appeal.StateWaitReviewer,
			PaperId:    paperId,
			ExaminerId: p.ExaminerId,
			ReviewerId: p.ReviewerId,
			TeacherId:  teacherId,
			AppealInfo: info,
		}
		var ap *types.ModelAppeal
		ap, err = appeal.AddAppeal(a)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		return ap, nil

	})

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
