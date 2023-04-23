package main

import (
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/config"
	"github.com/xihui-forever/mutualRead/impl"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/types"
	"reflect"
	"strconv"
)

func main() {
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

	err = role.Load()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	app := goon.New()

	app.PreUse("/", func(ctx *goon.Ctx) error {
		ctx.SetResHeader("Origin", "*")
		ctx.SetResHeader("Access-Control-Allow-Origin", "*")
		ctx.SetResHeader("Access-Control-Allow-Credentials", "true")
		ctx.SetResHeader("Access-Control-Expose-Headers", "")
		ctx.SetResHeader("Access-Control-Allow-Methods", "*")
		ctx.SetResHeader("Access-Control-Allow-Headers", "*")
		if ctx.Method() == "OPTIONS" {
			ctx.SetStatusCode(200)
			return ctx.Send("")
		}
		return ctx.Next()
	})

	app.PreUse("/", func(ctx *goon.Ctx) error {
		path := ctx.Path()
		flag, err := role.CheckPermission(role.RoleTypePublic, path)
		if flag {
			return ctx.Next()
		}
		/*if err != role.ErrRolePermExists {
			log.Errorf("err:%v", err)
			return err
		}*/

		var loginReq impl.LoginSession
		sid := ctx.GetReqHeader("X-Session-Id")
		s, err := session.GetSession(sid)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = sonic.UnmarshalString(s, &loginReq)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		ctx.SetReqHeader("X-Role-Type", strconv.Itoa(loginReq.RoleType))

		roleType := loginReq.RoleType
		log.Info(roleType)
		_, err = role.CheckPermission(roleType, path)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return ctx.Next()
	})

	for _, value := range impl.CmdList {
		app.Post(value.Path, AddPost(value.Logic))
	}

	err = app.ListenAndServe(":8080")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
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
		log.Error(x.PkgPath())
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
		req := reflect.New(in)
		err := ctx.ParseBody(req.Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// 调用逻辑
		ret := lv.Call([]reflect.Value{reflect.ValueOf(ctx), req})
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
