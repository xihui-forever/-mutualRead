package impl

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/exam"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/types"
	"strconv"
)

func init() {
	rpc.Register(types.CmdPathAddPaper, AddPaper, types.RoleTypeTeacher)
	rpc.Register(types.CmdPathListPaperTeacher, ListPaperTeacher, types.RoleTypeTeacher)
	rpc.Register(types.CmdPathListPaperExaminer, ListPaperExaminer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathGetPaperTeacher, GetPaperTeacher, types.RoleTypeTeacher)
	rpc.Register(types.CmdPathGetPaperExaminer, GetPaperExaminer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathGetPaperReviewer, GetPaperReviewer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathDelPaperTeacher, DelPaperTeacher, types.RoleTypeTeacher)
}

func AddPaper(ctx *goon.Ctx, req *types.AddPaperReq) (*types.AddPaperRsp, error) {
	var rsp types.AddPaperRsp

	e, err := exam.Get(req.Paper.ExamId)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if e.TeacherId != ctx.GetUint64(types.HeaderUserId) {
		return nil, exam.ErrExamNotExist
	}

	req.Paper.TeacherId = e.Id

	rsp.Paper, err = paper.Add(req.Paper)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// TODO 通知

	return &rsp, nil
}

func ListPaperTeacher(ctx *goon.Ctx, req *types.ListPaperTeacherReq) (*types.ListPaperTeacherRsp, error) {
	var rsp types.ListPaperTeacherRsp

	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListPaper_OptionTeacherId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var err error
	rsp.PaperList, rsp.Page, err = paper.ListPaper(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func ListPaperExaminer(ctx *goon.Ctx, req *types.ListPaperExaminerReq) (*types.ListPaperExaminerRsp, error) {
	var rsp types.ListPaperExaminerRsp

	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListPaper_OptionTeacherId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var err error
	rsp.PaperList, rsp.Page, err = paper.ListPaper(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func GetPaperTeacher(ctx *goon.Ctx, req *types.GetPaperTeacherReq) (*types.GetPaperTeacherRsp, error) {
	var rsp types.GetPaperTeacherRsp

	opt := &types.ListOption{
		Options: []types.Option{
			{
				Key: types.ListPaper_OptionId,
				Val: strconv.FormatUint(req.Id, 10),
			},
			{
				Key: types.ListPaper_OptionTeacherId,
				Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
			},
		},
		Limit: 1,
	}

	papers, _, err := paper.ListPaper(opt)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(papers) == 0 {
		return nil, paper.ErrPaperNotExist
	}

	rsp.Paper = papers[0]
	return &rsp, nil
}

func GetPaperExaminer(ctx *goon.Ctx, req *types.GetPaperExaminerReq) (*types.GetPaperExaminerRsp, error) {
	var rsp types.GetPaperExaminerRsp

	opt := &types.ListOption{
		Options: []types.Option{
			{
				Key: types.ListPaper_OptionId,
				Val: strconv.FormatUint(req.Id, 10),
			},
			{
				Key: types.ListPaper_OptionExaminerId,
				Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
			},
		},
		Limit: 1,
	}

	papers, _, err := paper.ListPaper(opt)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(papers) == 0 {
		return nil, paper.ErrPaperNotExist
	}

	rsp.Paper = papers[0]

	return &rsp, nil
}

func GetPaperReviewer(ctx *goon.Ctx, req *types.GetPaperReviewerReq) (*types.GetPaperReviewerRsp, error) {
	var rsp types.GetPaperReviewerRsp

	opt := &types.ListOption{
		Options: []types.Option{
			{
				Key: types.ListPaper_OptionId,
				Val: strconv.FormatUint(req.Id, 10),
			},
			{
				Key: types.ListPaper_OptionReviewerId,
				Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
			},
		},
		Limit: 1,
	}

	papers, _, err := paper.ListPaper(opt)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(papers) == 0 {
		return nil, paper.ErrPaperNotExist
	}

	rsp.Paper = papers[0]

	return &rsp, nil
}

func DelPaperTeacher(ctx *goon.Ctx, req *types.DelPaperTeacherReq) error {
	return paper.DelPaper(req.Id, ctx.GetUint64(types.HeaderUserId))
}
