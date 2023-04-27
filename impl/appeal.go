package impl

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/appeal"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/types"
	"strconv"
)

func init() {
	rpc.Register(types.CmdPathAddAppeal, AddAppeal, types.RoleTypeStudent)
	rpc.Register(types.CmdPathListAppealExaminer, ListAppealExaminer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathListAppealReviewer, ListAppealReviewer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathListAppealTeacher, ListAppealTeacher, types.RoleTypeTeacher)

	rpc.Register(types.CmdPathGetAppealExaminer, GetAppealExaminer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathGetAppealReviewer, GetAppealReviewer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathGetAppealTeacher, GetAppealTeacher, types.RoleTypeTeacher)
}

func AddAppeal(ctx *goon.Ctx, req *types.AddAppealReq) (*types.AddAppealRsp, error) {
	var rsp types.AddAppealRsp

	p, err := paper.Get(req.PaperId)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if p.ExaminerId != ctx.GetUint64(types.HeaderUserId) {
		return nil, paper.ErrPaperNotExist
	}

	a, err := appeal.AddAppeal(&types.ModelAppeal{
		State:      types.AppealStateWaitReviewer,
		PaperId:    p.Id,
		ExaminerId: p.ExaminerId,
		ReviewerId: p.ReviewerId,
		TeacherId:  p.TeacherId,
		AppealInfo: req.AppealInfo,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Appeal = a
	return &rsp, nil
}

func ListAppealExaminer(ctx *goon.Ctx, req *types.ListAppealExaminerReq) (*types.ListAppealExaminerRsp, error) {
	var rsp types.ListAppealExaminerRsp
	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListPaper_OptionExaminerId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var err error
	rsp.List, rsp.Page, err = appeal.ListAppeal(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if req.ShowPaper && len(rsp.List) > 0 {
		var papers []*types.ModelPaper
		err = db.Where("id in (?)", utils.PluckUint64(rsp.List, "PaperId")).Find(&papers).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.PaperMap = make(map[uint64]*types.ModelPaper, len(papers))
		for _, a := range papers {
			rsp.PaperMap[a.Id] = a
		}
	}

	return &rsp, nil
}

func ListAppealReviewer(ctx *goon.Ctx, req *types.ListAppealReviewerReq) (*types.ListAppealReviewerRsp, error) {
	var rsp types.ListAppealReviewerRsp
	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListPaper_OptionReviewerId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var err error
	rsp.List, rsp.Page, err = appeal.ListAppeal(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if req.ShowPaper && len(rsp.List) > 0 {
		var papers []*types.ModelPaper
		err = db.Where("id in (?)", utils.PluckUint64(rsp.List, "PaperId")).Find(&papers).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.PaperMap = make(map[uint64]*types.ModelPaper, len(papers))
		for _, a := range papers {
			rsp.PaperMap[a.Id] = a
		}
	}

	return &rsp, nil
}

func ListAppealTeacher(ctx *goon.Ctx, req *types.ListAppealTeacherReq) (*types.ListAppealTeacherRsp, error) {
	var rsp types.ListAppealTeacherRsp
	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListPaper_OptionTeacherId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var err error
	rsp.List, rsp.Page, err = appeal.ListAppeal(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if req.ShowPaper && len(rsp.List) > 0 {
		var papers []*types.ModelPaper
		err = db.Where("id in (?)", utils.PluckUint64(rsp.List, "PaperId")).Find(&papers).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.PaperMap = make(map[uint64]*types.ModelPaper, len(papers))
		for _, a := range papers {
			rsp.PaperMap[a.Id] = a
		}
	}

	return &rsp, nil
}

func GetAppealTeacher(ctx *goon.Ctx, req *types.GetAppealTeacherReq) (*types.GetAppealTeacherRsp, error) {
	var rsp types.GetAppealTeacherRsp

	list, _, err := appeal.ListAppeal(&types.ListOption{
		Options: []types.Option{
			{
				Key: types.ListAppeal_OptionId,
				Val: strconv.FormatUint(req.Id, 10),
			},
			{
				Key: types.ListAppeal_OptionTeacherId,
				Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
			},
		},
		Limit: 1,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(list) == 0 {
		return nil, appeal.ErrAppealNotExist
	}

	rsp.Appeal = list[0]

	if req.ShowPaper {
		rsp.Paper, err = paper.Get(rsp.Appeal.PaperId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}

	return &rsp, nil
}

func GetAppealReviewer(ctx *goon.Ctx, req *types.GetAppealReviewerReq) (*types.GetAppealReviewerRsp, error) {
	var rsp types.GetAppealReviewerRsp

	list, _, err := appeal.ListAppeal(&types.ListOption{
		Options: []types.Option{
			{
				Key: types.ListAppeal_OptionId,
				Val: strconv.FormatUint(req.Id, 10),
			},
			{
				Key: types.ListAppeal_OptionReviewerId,
				Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
			},
		},
		Limit: 1,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(list) == 0 {
		return nil, appeal.ErrAppealNotExist
	}

	rsp.Appeal = list[0]

	if req.ShowPaper {
		rsp.Paper, err = paper.Get(rsp.Appeal.PaperId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}

	return &rsp, nil
}

func GetAppealExaminer(ctx *goon.Ctx, req *types.GetAppealExaminerReq) (*types.GetAppealExaminerRsp, error) {
	var rsp types.GetAppealExaminerRsp

	list, _, err := appeal.ListAppeal(&types.ListOption{
		Options: []types.Option{
			{
				Key: types.ListAppeal_OptionId,
				Val: strconv.FormatUint(req.Id, 10),
			},
			{
				Key: types.ListAppeal_OptionExaminerId,
				Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
			},
		},
		Limit: 1,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if len(list) == 0 {
		return nil, appeal.ErrAppealNotExist
	}

	rsp.Appeal = list[0]

	if req.ShowPaper {
		rsp.Paper, err = paper.Get(rsp.Appeal.PaperId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}

	return &rsp, nil
}
