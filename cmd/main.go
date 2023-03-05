package main

import (
	"reflect"

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
		app.Post(value.Path, AddPost(value.Logic))
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

func AddPost(logic interface{}) func(ctx *goon.Ctx) error {
	lv := reflect.ValueOf(logic)
	lt := reflect.TypeOf(logic)

	if lt.Kind() != reflect.Func {
		panic("parameter is not func")
	}

	if lt.NumIn() != 2 {
		panic("num of func in can not be empty")
	}

	if lt.NumOut() != 2 {
		panic("num of func out can not be empty")
	}

	// 第一个入参必须是*goon.Ctx
	x := lt.In(0)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Kind() != reflect.Struct {
		panic("first in is must *Ctx")
	}

	if x.PkgPath() != "github.com/xihui-forever/goon" {
		panic("first in is must *goon.Ctx")
	}

	if x.Name() != "Ctx" {
		panic("first in is must *goon.Ctx")
	}

	// 第二个入参是struct
	in := lt.In(1)
	for in.Kind() == reflect.Ptr {
		in = in.Elem()
	}

	if in.Kind() != reflect.Struct {
		panic("second in is must struct")
	}

	// 第一个出参必须是struct
	out := lt.Out(0)
	for out.Kind() == reflect.Ptr {
		out = out.Elem()
	}

	if out.Kind() != reflect.Struct {
		panic("first out is must struct")
	}

	// 第二个出参必须是error
	x = lt.Out(1)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Name() != "error" {
		panic("second out is must error")
	}

	return func(ctx *goon.Ctx) error {
		req := reflect.New(out)
		err := ctx.ParseBody(req.Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// 调用逻辑
		ret := lv.Call([]reflect.Value{reflect.ValueOf(ctx), req.Elem()})
		if !ret[1].IsNil() {
			log.Errorf("err:%v", ret[1].Interface())

			if x, ok := ret[1].Interface().(*types.Error); ok {
				return ctx.Json(x)
			} else {
				return ctx.Json(&types.Error{
					Code: types.SysError,
					Msg:  ret[1].Interface().(error).Error(),
				})
			}
		}
		// 返回数据

		return ctx.Json(map[string]any{
			"code": 0,
			"msg":  "ok",
			"data": ret[0].Interface(),
		})
	}
}
