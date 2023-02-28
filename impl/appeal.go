package impl

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/appeal"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/types"
)

type AddAppealReq struct {
	PaperId   uint64 `json:"paperId"`
	TeacherId uint64 `json:"teacherId"`
	Info      string `json:"info"`
}

func init() {
	CmdList = append(CmdList, Cmd{
		Path:  "/appeal_add",
		Role:  2,
		Logic: AddAppeal,
	})
}

func AddAppeal(ctx *goon.Ctx, req AddAppealReq) (*types.ModelAppeal, error) {
	err := ctx.ParseBody(req)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	var p *types.ModelPaper
	p, err = paper.GetPaper(req.PaperId, req.TeacherId)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	a := types.ModelAppeal{
		State:      appeal.StateWaitReviewer,
		PaperId:    req.PaperId,
		ExaminerId: p.ExaminerId,
		ReviewerId: p.ReviewerId,
		TeacherId:  req.TeacherId,
		AppealInfo: req.Info,
	}
	var ap *types.ModelAppeal
	ap, err = appeal.AddAppeal(a)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return ap, nil

}
