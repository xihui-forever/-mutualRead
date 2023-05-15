package rpc

import (
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/public"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/types"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func Load() {
	goon.PreUse("/", func(ctx *goon.Ctx) error {
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

	// 前端的静态文件
	goon.Use("/", func(ctx *goon.Ctx) error {
		path := ctx.Path()
		if path == "/" {
			path = "/index.html"
		}

		// 判断一下静态文件是否存在
		buf, err := public.Public.ReadFile(filepath.ToSlash(filepath.Join("build", path)))
		if err != nil {
			log.Errorf("err:%v", err)
			return ctx.Next()
		}

		ctx.SetResHeader("Cache-Control", "public, max-age=86400")

		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			var buffer []byte
			if len(buffer) < 512 {
				buffer = buf
			} else {
				buffer = buf[:512]
			}
			contentType = http.DetectContentType(buffer)
		}

		ctx.SetResHeader("Content-Type", contentType)

		return ctx.Write(buf)
	})

	// 后端的校验逻辑
	goon.PreUse("/", func(ctx *goon.Ctx) error {
		path := ctx.Path()

		if strings.HasPrefix(path, types.CmdPathResourceGet) {
			return ctx.Next()
		}

		flag, err := role.CheckPermission(types.RoleTypePublic, path)
		if flag {
			return ctx.Next()
		}

		sessionBuf, err := session.GetSession(ctx.GetReqHeader("X-Session-Id"))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		var sess types.LoginSession
		err = sonic.UnmarshalString(sessionBuf, &sess)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		ctx.SetReqHeader(types.HeaderUserId, strconv.FormatUint(sess.Id, 10))
		ctx.SetReqHeader(types.HeaderRoleType, strconv.Itoa(sess.RoleType))

		ctx.Set(types.HeaderUserId, sess.Id)
		ctx.Set(types.HeaderRoleType, sess.RoleType)

		roleType := sess.RoleType
		_, err = role.CheckPermission(roleType, path)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return ctx.Next()
	})

	for _, cmd := range CmdList {
		goon.Register(cmd.Method, cmd.Path, GenHandler(cmd.Logic))

		for _, r := range cmd.Roles {
			_, _ = role.BatchAddRolePerm(r, []string{cmd.Path})
		}
	}
}

func HandleError(ctx *goon.Ctx, err error) error {
	if err == nil {
		return ctx.Json(&types.Error{
			TranceId: log.GetTrace(),
		})
	}

	if x, ok := err.(*types.Error); ok {
		x.TranceId = log.GetTrace()
		return ctx.Json(x)
	}

	return ctx.Json(&types.Error{
		Code:     types.SysError,
		Msg:      err.Error(),
		TranceId: log.GetTrace(),
	})
}

func GenHandler(logic any) goon.Handler {
	switch logic.(type) {
	case goon.Handler:
		return func(ctx *goon.Ctx) error {
			return HandleError(ctx, logic.(goon.Handler)(ctx))
		}
	default:
		lt := reflect.TypeOf(logic)
		lv := reflect.ValueOf(logic)
		if lt.Kind() != reflect.Func {
			panic("parameter is not func")
		}

		// 不管怎么样，第一个参数都是一定要存在的
		x := lt.In(0)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() != reflect.Struct {
			panic("first in is must *github.com/xihui-forever/goon.Ctx")
		}
		if x.Name() != "Ctx" {
			panic("first in is must *github.com/xihui-forever/goon.Ctx")
		}
		if x.PkgPath() != "github.com/xihui-forever/goon" {
			panic("first in is must *github.com/xihui-forever/goon.Ctx")
		}

		// 按照三种不同的情况来处理
		if lt.NumIn() == 1 && lt.NumOut() == 1 {
			// 一个入参，一个出参，那就是Handler本身了
			// 但是这里还要判断一下，出参是否是error

			x = lt.Out(0)
			for x.Kind() == reflect.Ptr {
				x = x.Elem()
			}

			if x.Name() != "error" {
				panic("out is must error")
			}

			return func(ctx *goon.Ctx) error {
				out := lv.Call([]reflect.Value{reflect.ValueOf(ctx)})
				if out[0].IsNil() {
					if ctx.ContentLen() > 0 {
						return nil
					}

					return ctx.Json(&types.Error{
						TranceId: log.GetTrace(),
					})
				}

				return HandleError(ctx, out[0].Interface().(error))
			}
		} else if lt.NumIn() == 2 && lt.NumOut() == 1 {
			// 两个入参，一个出参，那就是需要解析请求参数的

			// 先判断出参是否是error
			x = lt.Out(0)
			for x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			if x.Name() != "error" {
				panic("out is must error")
			}

			// 处理一下入参
			x = lt.In(1)
			for x.Kind() == reflect.Ptr {
				x = x.Elem()
			}

			if x.Kind() != reflect.Struct {
				panic("2rd in is must struct")
			}

			return func(ctx *goon.Ctx) error {
				req := reflect.New(x)
				err := ctx.ParseBody(req.Interface())
				if err != nil {
					return err
				}

				err = utils.Validate(req.Interface())
				if err != nil {
					log.Errorf("err:%v", err)
					return ctx.Json(&types.Error{
						Code:     types.ErrInvalidParam,
						Msg:      err.Error(),
						TranceId: log.GetTrace(),
					})
				}

				out := lv.Call([]reflect.Value{reflect.ValueOf(ctx), req})
				if out[0].IsNil() {
					return ctx.Json(&types.Error{
						TranceId: log.GetTrace(),
					})
				}

				return HandleError(ctx, out[0].Interface().(error))
			}
		} else if lt.NumIn() == 1 && lt.NumOut() == 2 {
			// 一个入参，两个出参，那就是需要返回数据的

			// 先判断第二个出参是否是error
			x = lt.Out(1)
			for x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			if x.Name() != "error" {
				panic("out 1 is must error")
			}

			// 处理一下返回
			x = lt.Out(0)
			for x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			if x.Kind() != reflect.Struct {
				panic("out 0 is must struct")
			}

			return func(ctx *goon.Ctx) error {
				out := lv.Call([]reflect.Value{reflect.ValueOf(ctx)})
				if out[1].IsNil() {
					return ctx.JsonWithPerm(ctx.GetReqHeader(types.HeaderRoleType), &types.Error{
						TranceId: log.GetTrace(),
						Data:     out[0].Interface(),
					})
				}

				return HandleError(ctx, out[1].Interface().(error))
			}
		} else if lt.NumIn() == 2 && lt.NumOut() == 2 {
			// 两个入参，两个出参，那就是需要解析请求参数，返回数据的

			// 先判断第二个出参是否是error
			x = lt.Out(1)
			for x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			if x.Name() != "error" {
				panic("out 1 is must error")
			}

			// 处理一下入参
			in := lt.In(1)
			for in.Kind() == reflect.Ptr {
				in = in.Elem()
			}
			if in.Kind() != reflect.Struct {
				panic("2rd in is must struct")
			}

			// 处理一下返回
			out := lt.Out(0)
			for out.Kind() == reflect.Ptr {
				out = out.Elem()
			}
			if out.Kind() != reflect.Struct {
				panic("out 0 is must struct")
			}

			return func(ctx *goon.Ctx) error {
				req := reflect.New(in)
				err := ctx.ParseBody(req.Interface())
				if err != nil {
					return err
				}

				err = utils.Validate(req.Interface())
				if err != nil {
					log.Errorf("err:%v", err)
					return ctx.Json(&types.Error{
						Code:     types.ErrInvalidParam,
						Msg:      err.Error(),
						TranceId: log.GetTrace(),
					})
				}

				out := lv.Call([]reflect.Value{reflect.ValueOf(ctx), req})
				if out[1].IsNil() {
					return ctx.JsonWithPerm(ctx.GetReqHeader(types.HeaderRoleType), &types.Error{
						TranceId: log.GetTrace(),
						Data:     out[0].Interface(),
					})
				}

				return HandleError(ctx, out[1].Interface().(error))
			}
		} else {
			panic("func is not support")
		}
	}
}

type Cmd struct {
	Path   string
	Roles  []int
	Logic  interface{} // func(ctx, req) (resp, err)
	Method goon.Method
}

var CmdList = []Cmd{}

func Post(path string, logic any, roles ...int) {
	CmdList = append(CmdList, Cmd{
		Method: goon.MethodPost,
		Path:   path,
		Roles:  roles,
		Logic:  logic,
	})
}

func Put(path string, logic any, roles ...int) {
	CmdList = append(CmdList, Cmd{
		Method: goon.MethodPut,
		Path:   path,
		Roles:  roles,
		Logic:  logic,
	})
}

func Get(path string, logic any, roles ...int) {
	CmdList = append(CmdList, Cmd{
		Method: goon.MethodGet,
		Path:   path,
		Roles:  roles,
		Logic:  logic,
	})
}

func Use(path string, logic any, roles ...int) {
	CmdList = append(CmdList, Cmd{
		Method: goon.MethodUse,
		Path:   path,
		Roles:  roles,
		Logic:  logic,
	})
}
