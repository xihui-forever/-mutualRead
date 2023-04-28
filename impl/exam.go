package impl

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/exam"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/types"
	"strconv"
)

func init() {
	rpc.Register(types.CmdPathAddExam, AddExam, types.RoleTypeTeacher)
	rpc.Register(types.CmdPathSetExam, SetExam, types.RoleTypeTeacher)
	rpc.Register(types.CmdPathDelExam, DelExam, types.RoleTypeTeacher)
	rpc.Register(types.CmdPathListExam, ListExam, types.RoleTypeTeacher)
}

func AddExam(ctx *goon.Ctx, req *types.AddExamReq) (*types.ModelExam, error) {
	req.Exam.TeacherId = ctx.GetUint64(types.HeaderUserId)
	exam, err := exam.Add(req.Exam)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return exam, nil
}

func SetExam(ctx *goon.Ctx, req *types.SetExamReq) (*types.ModelExam, error) {
	// 存在权限问题，所以要先看下用户是否对的齐
	teacherId := ctx.GetUint64(types.HeaderUserId)

	e, err := exam.Get(req.Exam.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if e.TeacherId != teacherId {
		return nil, exam.ErrExamNotExist
	}

	req.Exam.TeacherId = teacherId

	e, err = exam.Set(req.Exam)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return e, nil
}

func DelExam(ctx *goon.Ctx, req *types.ModelExam) error {
	// 存在权限问题，所以要先看下用户是否对的齐
	teacherId := ctx.GetUint64(types.HeaderUserId)

	e, err := exam.Get(req.Id)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if e.TeacherId != teacherId {
		return exam.ErrExamNotExist
	}

	err = exam.Del(req.Id)
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

func ListExam(ctx *goon.Ctx, req *types.ListExamReq) (*types.ListExamRsp, error) {
	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListExam_OptionTeacherId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var rsp types.ListExamRsp
	var err error
	rsp.Exams, rsp.Page, err = exam.List(req.Options)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	return &rsp, nil
}
