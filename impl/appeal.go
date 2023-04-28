package impl

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/darabuchi/utils/mq"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/appeal"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/types"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func init() {
	rpc.Register(types.CmdPathAddAppeal, AddAppeal, types.RoleTypeStudent)
	rpc.Register(types.CmdPathListAppealExaminer, ListAppealExaminer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathListAppealReviewer, ListAppealReviewer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathListAppealTeacher, ListAppealTeacher, types.RoleTypeTeacher)

	rpc.Register(types.CmdPathGetAppealExaminer, GetAppealExaminer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathGetAppealReviewer, GetAppealReviewer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathGetAppealTeacher, GetAppealTeacher, types.RoleTypeTeacher)

	rpc.Register(types.CmdPathSetAppealExaminer, SetAppealExaminer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathSetAppealReviewer, SetAppealReviewer, types.RoleTypeStudent)
	rpc.Register(types.CmdPathSetAppealTeacher, SetAppealTeacher, types.RoleTypeTeacher)
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

	_, err = mq.Publish(types.EventAppealChangedTopic, types.EventAppealChanged{
		Appeal: a,
	})
	if err != nil {
		log.Errorf("err:%v", err)
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

func SetAppealExaminer(ctx *goon.Ctx, req *types.SetAppealExaminerReq) error {
	a, err := appeal.Get(req.AppealId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a.ExaminerId != ctx.GetUint64(types.HeaderUserId) {
		return appeal.ErrAppealNotExist
	}

	if a.State != types.AppealStateWaitReviewer {
		return appeal.ErrAppealAlreadyHanded
	}

	res := db.Model(&types.ModelAppeal{}).
		Where("id = ?", req.AppealId).
		Where("state = ?", types.AppealStateWaitReviewer).
		Updates(map[string]any{
			"change_at":   time.Now().Unix(),
			"appeal_info": req.AppealInfo,
		})
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}

	if res.RowsAffected == 0 {
		return appeal.ErrAppealAlreadyHanded
	}

	return nil
}

func SetAppealReviewer(ctx *goon.Ctx, req *types.SetAppealReviewerReq) error {
	a, err := appeal.Get(req.AppealId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a.ReviewerId != ctx.GetUint64(types.HeaderUserId) {
		return appeal.ErrAppealNotExist
	}

	if a.State != types.AppealStateWaitReviewer {
		return appeal.ErrAppealAlreadyHanded
	}

	res := db.Model(&types.ModelAppeal{}).
		Where("id = ?", req.AppealId).
		Where("state = ?", types.AppealStateWaitReviewer).
		Updates(map[string]any{
			"state": types.AppealStateWaitTeacher,

			"reviewer_at": time.Now().Unix(),
			"review_info": req.ReviewInfo,
		})
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}

	if res.RowsAffected == 0 {
		return appeal.ErrAppealAlreadyHanded
	}

	a, err = appeal.Get(req.AppealId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = mq.Publish(types.EventAppealChangedTopic, types.EventAppealChanged{
		Appeal: a,
	})
	if err != nil {
		log.Errorf("err:%v", err)
	}

	return nil
}

func SetAppealTeacher(ctx *goon.Ctx, req *types.SetAppealTeacherReq) error {
	a, err := appeal.Get(req.AppealId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a.TeacherId != ctx.GetUint64(types.HeaderUserId) {
		return appeal.ErrAppealNotExist
	}

	if a.State != types.AppealStateWaitTeacher {
		return appeal.ErrAppealAlreadyHanded
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&types.ModelAppeal{}).
			Where("id = ?", req.AppealId).
			Where("state = ?", types.AppealStateWaitTeacher).
			Updates(map[string]any{
				"state": types.AppealStateFinish,

				"result_at":     time.Now().Unix(),
				"appeal_result": req.AppealResult,
				"grade":         req.Grade,
			})
		if res.Error != nil {
			log.Errorf("err:%v", res.Error)
			return res.Error
		}

		if res.RowsAffected == 0 {
			return appeal.ErrAppealAlreadyHanded
		}

		err = tx.Model(&types.ModelPaper{}).
			Where("id = ?", a.PaperId).
			Update("grade", gorm.Expr("grade += ?", req.Grade)).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	a, err = appeal.Get(req.AppealId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = mq.Publish(types.EventAppealChangedTopic, types.EventAppealChanged{
		Appeal: a,
	})
	if err != nil {
		log.Errorf("err:%v", err)
	}

	return nil
}

func RecallAppeal(ctx *goon.Ctx, req *types.RecallAppealReq) error {
	a, err := appeal.Get(req.AppealId)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a.ExaminerId != ctx.GetUint64(types.HeaderUserId) {
		return appeal.ErrAppealNotExist
	}

	if a.State != types.AppealStateWaitReviewer {
		return appeal.ErrAppealAlreadyHanded
	}

	res := db.Model(&types.ModelAppeal{}).
		Where("id = ?", req.AppealId).
		Where("state = ?", types.AppealStateWaitReviewer).
		Updates(map[string]any{
			"state": types.AppealStateRecall,
		})
	if res.Error != nil {
		log.Errorf("err:%v", res.Error)
		return res.Error
	}

	if res.RowsAffected == 0 {
		return appeal.ErrAppealAlreadyHanded
	}

	return nil
}
