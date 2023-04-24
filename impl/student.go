package impl

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/student"
	"github.com/xihui-forever/mutualRead/types"
)

func init() {
	rpc.Register(types.CmdPathGetStudent, GetStudent, role.RoleTypeStudent)
	rpc.Register(types.CmdPathChangeStudent, ChangeStudent, role.RoleTypeStudent)

	rpc.Register(types.CmdPathAddStudentAdmin, AddStudentAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathGetStudentAdmin, GetStudentAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathSetStudentAdmin, SetStudentAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathDelStudentAdmin, DeleteStudentAdmin, role.RoleTypeAdmin)
	rpc.Register(types.CmdPathListStudentAdmin, ListStudentAdmin, role.RoleTypeAdmin)
}

func GetStudent(ctx *goon.Ctx) (*types.GetStudentRsp, error) {
	student, err := student.Get(ctx.GetWithDef(types.HeaderUserId, 0).(uint64))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &types.GetStudentRsp{
		Student: student,
	}, nil
}

func ChangeStudent(ctx *goon.Ctx, req *types.ChangeStudentReq) error {
	switch req.ChangeType {
	case types.StudentChangeTypeEmail:
		err := student.ChangeEmail(ctx.GetWithDef(types.HeaderUserId, 0).(uint64), req.Email)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	default:
		log.Errorf("invalid change type:%v", req.ChangeType)
		return types.CreateErrorWithMsg(types.ErrInvalidParam, "invalid change type")
	}

	return nil
}

func AddStudentAdmin(ctx *goon.Ctx, req *types.AddStudentAdminReq) (*types.AddStudentAdminRsp, error) {
	student, err := student.AddStudent(*req.Student)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &types.AddStudentAdminRsp{
		Student: student,
	}, nil
}

func GetStudentAdmin(ctx *goon.Ctx, req *types.GetStudentAdminReq) (*types.GetStudentAdminRsp, error) {
	student, err := student.Get(req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &types.GetStudentAdminRsp{
		Student: student,
	}, nil
}

func SetStudentAdmin(ctx *goon.Ctx, req *types.SetStudentAdminReq) (*types.SetStudentAdminRsp, error) {
	student, err := student.Set(req.Student)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &types.SetStudentAdminRsp{
		Student: student,
	}, nil
}

func DeleteStudentAdmin(ctx *goon.Ctx, req *types.DelStudentAdminReq) error {
	err := student.Del(req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func ListStudentAdmin(ctx *goon.Ctx, req *types.ListStudentAdminReq) (*types.ListStudentAdminRsp, error) {
	var rsp types.ListStudentAdminRsp
	var err error
	rsp.Students, rsp.Page, err = student.List(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
