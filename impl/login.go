package impl

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/xtime"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/admin"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/student"
	"github.com/xihui-forever/mutualRead/teacher"
	"github.com/xihui-forever/mutualRead/types"
)

type LoginBaseRsp interface {
	GetId() uint64
}

var (
	LoginHandlerMap  = map[int]func(username string, password string) (LoginBaseRsp, error){}
	ResetPasswordMap = map[int]func(username, password string) error{}
)

func init() {
	LoginHandlerMap[types.RoleTypeTeacher] = func(username string, password string) (LoginBaseRsp, error) {
		data, err := teacher.GetTeacher(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = teacher.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return data, nil
	}
	ResetPasswordMap[types.RoleTypeAdmin] = admin.ResetPassword

	LoginHandlerMap[types.RoleTypeAdmin] = func(username string, password string) (LoginBaseRsp, error) {
		data, err := admin.Get(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = admin.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return data, nil
	}
	ResetPasswordMap[types.RoleTypeTeacher] = teacher.ResetPassword

	LoginHandlerMap[types.RoleTypeStudent] = func(username string, password string) (LoginBaseRsp, error) {
		data, err := student.GetStudent(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = student.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return data, nil
	}
	ResetPasswordMap[types.RoleTypeStudent] = student.ResetPassword

	rpc.Post("/login", LoginHandler, types.RoleTypePublic)
	rpc.Post("/reset/password", ResetPassword, types.RoleTypeStudent, types.RoleTypeTeacher, types.RoleTypeAdmin)
	rpc.Post("/change-password", ChangePassword, types.RoleTypeStudent, types.RoleTypeTeacher, types.RoleTypeAdmin)
}

func LoginHandler(ctx *goon.Ctx, req *types.LoginReq) (*types.LoginRsp, error) {
	var resp types.LoginRsp
	logic, ok := LoginHandlerMap[req.RoleType]
	if !ok {
		log.Errorf("login type %v not found", req.RoleType)
		return nil, types.CreateError(types.ErrLoginTypeNotFound)
	}

	info, err := logic(req.Username, req.Password)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Token, err = session.GenSession(&types.LoginSession{
		RoleType: req.RoleType,
		Id:       info.GetId(),
	}, xtime.Day)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Expire = uint32(xtime.Day.Seconds())
	resp.Info = info
	return &resp, nil
}

func ResetPassword(ctx *goon.Ctx, req *types.ResetPasswordReq) error {
	switch ctx.GetInt(types.HeaderRoleType) {
	case types.RoleTypeAdmin:
		logic, ok := ResetPasswordMap[req.RoleType]
		if !ok {
			log.Errorf("login type %v not found", req.RoleType)
			return types.CreateError(types.ErrLoginTypeNotFound)
		}

		err := logic(req.Username, req.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	case types.RoleTypeStudent:
		s, err := student.Get(ctx.GetUint64(types.HeaderUserId))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return student.ResetPassword(s.StudentId, req.Password)
	case types.RoleTypeTeacher:
		s, err := teacher.Get(ctx.GetUint64(types.HeaderUserId))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return teacher.ResetPassword(s.TeacherId, req.Password)
	default:
		ctx.SetStatusCode(403)
		return types.CreateErrorWithMsg(-1, "user not have perm")
	}

	return nil
}

func ChangePassword(ctx *goon.Ctx, req *types.ChangePasswordReq) error {
	switch ctx.GetInt(types.HeaderRoleType) {
	case types.RoleTypeAdmin:
		return admin.ChangePassword(ctx.GetUint64(types.HeaderUserId), req.OldPassword, req.NewPassword)
	case types.RoleTypeStudent:
		return student.ChangePassword(ctx.GetUint64(types.HeaderUserId), req.OldPassword, req.NewPassword)
	case types.RoleTypeTeacher:
		return teacher.ChangePassword(ctx.GetUint64(types.HeaderUserId), req.OldPassword, req.NewPassword)
	default:
		ctx.SetStatusCode(403)
		return types.CreateErrorWithMsg(-1, "user not have perm")
	}
}
