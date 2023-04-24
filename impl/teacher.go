package impl

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/teacher"
	"github.com/xihui-forever/mutualRead/types"
)

func init() {
	rpc.Register(types.CmdPathGetTeacher, GetTeacher, role.RoleTypeTeacher)
	rpc.Register(types.CmdPathChangeTeacher, ChangeTeacher, role.RoleTypeTeacher)

	rpc.Register(types.CmdPathGetTeacherAdmin, GetTeacherAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathAddTeacherAdmin, AddTeacherAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathSetTeacherAdmin, SetTeacherAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathDelTeacherAdmin, DelTeacherAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathListTeacherAdmin, ListTeacherAdmin, role.RoleTypeAdmin)
}

func GetTeacher(ctx *goon.Ctx) (*types.GetTeacherRsp, error) {
	var rsp types.GetTeacherRsp

	teacher, err := teacher.Get(ctx.GetWithDef(types.HeaderUserId, 0).(uint64))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Teacher = teacher
	return &rsp, nil
}

func ChangeTeacher(ctx *goon.Ctx, req *types.ChangeTeacherReq) error {
	switch req.ChangeType {
	case types.TeacherChangeTypeEmail:
		err := teacher.ChangeEmail(ctx.GetWithDef(types.HeaderUserId, 0).(uint64), req.Email)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	default:
		return types.CreateErrorWithMsg(types.ErrInvalidParam, "invalid change type")
	}

	return nil
}

func GetTeacherAdmin(ctx *goon.Ctx, req *types.GetTeacherAdminReq) (*types.GetTeacherAdminRsp, error) {
	var rsp types.GetTeacherAdminRsp

	teacher, err := teacher.Get(req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Teacher = teacher
	return &rsp, nil
}

func AddTeacherAdmin(ctx *goon.Ctx, req *types.AddTeacherAdminReq) (*types.AddTeacherAdminRsp, error) {
	var rsp types.AddTeacherAdminRsp
	teacher, err := teacher.AddTeacher(req.Teacher)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Teacher = teacher
	return &rsp, nil
}

func SetTeacherAdmin(ctx *goon.Ctx, req *types.AddTeacherAdminReq) (*types.AddTeacherAdminRsp, error) {
	var rsp types.AddTeacherAdminRsp
	teacher, err := teacher.SetTeacher(req.Teacher)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Teacher = teacher
	return &rsp, nil
}

func DelTeacherAdmin(ctx *goon.Ctx, req *types.DelTeacherAdminReq) error {
	err := teacher.DelTeacher(req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func ListTeacherAdmin(ctx *goon.Ctx, req *types.ListTeacherAdminReq) (*types.ListTeacherAdminRsp, error) {
	var rsp types.ListTeacherAdminRsp

	var err error
	rsp.Teachers, rsp.Page, err = teacher.ListTeacher(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
