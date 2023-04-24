package impl

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/exam"
	"github.com/xihui-forever/mutualRead/types"
)

func init() {
	//CmdList = append(CmdList,
	//	Cmd{
	//		Path:  "/exam_list",
	//		Role:  2,
	//		Logic: GetExamList,
	//	},
	//	//Cmd{
	//	//	Path:  "/exam_get",
	//	//	Role:  2,
	//	//	Logic: GetTeacher,
	//	//},
	//	//Cmd{
	//	//	Path:  "/exam_add",
	//	//	Role:  2,
	//	//	Logic: ChangeTeacher,
	//	//},
	//	//Cmd{
	//	//	Path:  "/exam_del",
	//	//	Role:  2,
	//	//	Logic: GetExamList,
	//	//},
	//)
}

func AddExam(ctx *goon.Ctx, req *types.AddExamReq) (*types.ModelExam, error) {
	req.Exam.TeacherId = ctx.GetWithDef(types.HeaderUserId, 0).(uint64)
	exam, err := exam.Add(req.Exam)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return exam, nil
}

func SetExam(ctx *goon.Ctx, req types.SetExamReq) (*types.ModelExam, error) {
	// 存在权限问题，所以要先看下用户是否对的齐
	teacherId := ctx.GetWithDef(types.HeaderUserId, 0).(uint64)

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
	teacherId := ctx.GetWithDef(types.HeaderUserId, 0).(uint64)

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
		Val: ctx.GetReqHeader(types.HeaderUserId),
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
