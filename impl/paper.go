package impl

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/elliotchance/pie/v2"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/exam"
	"github.com/xihui-forever/mutualRead/paper"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/student"
	"github.com/xihui-forever/mutualRead/types"
	"strconv"
)

func init() {
	rpc.Post(types.CmdPathAddPaper, AddPaper, types.RoleTypeTeacher)
	rpc.Post(types.CmdPathListPaperTeacher, ListPaperTeacher, types.RoleTypeTeacher)
	rpc.Post(types.CmdPathListPaperExaminer, ListPaperExaminer, types.RoleTypeStudent)
	rpc.Post(types.CmdPathListPaperReviewer, ListPaperReviewer, types.RoleTypeStudent)
	rpc.Post(types.CmdPathGetPaperTeacher, GetPaperTeacher, types.RoleTypeTeacher)
	rpc.Post(types.CmdPathGetPaperExaminer, GetPaperExaminer, types.RoleTypeStudent)
	rpc.Post(types.CmdPathGetPaperReviewer, GetPaperReviewer, types.RoleTypeStudent)
	rpc.Post(types.CmdPathDelPaperTeacher, DelPaperTeacher, types.RoleTypeTeacher)
}

func AddPaper(ctx *goon.Ctx, req *types.AddPaperReq) (*types.AddPaperRsp, error) {
	var rsp types.AddPaperRsp

	e, err := exam.Get(req.ExamId)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if e.TeacherId != ctx.GetUint64(types.HeaderUserId) {
		return nil, exam.ErrExamNotExist
	}

	p := &types.ModelPaper{
		ExamId:    e.Id,
		TeacherId: e.TeacherId,
		Grade:     req.Grade,
		ImgUrl:    req.ImgUrl,
	}

	{
		s, err := student.GetByStrId(req.ExaminerId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		p.ExaminerId = s.Id
	}

	{
		s, err := student.GetByStrId(req.ReviewerId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		p.ReviewerId = s.Id
	}

	rsp.Paper, err = paper.Add(p)
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

	if req.ShowExam && len(rsp.PaperList) > 0 {
		var exams []*types.ModelExam
		err = db.Where("id in (?)", utils.PluckUint64(rsp.PaperList, "ExamId")).Find(&exams).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.ExamMap = make(map[uint64]*types.ModelExam, len(exams))
		for _, exam := range exams {
			rsp.ExamMap[exam.Id] = exam
		}
	}

	if req.ShowStudent && len(rsp.PaperList) > 0 {
		var studentIds []uint64
		studentIds = append(studentIds, utils.PluckUint64(rsp.PaperList, "ExaminerId")...)
		studentIds = append(studentIds, utils.PluckUint64(rsp.PaperList, "ReviewerId")...)

		studentIds = pie.Unique(studentIds)

		var students []*types.ModelStudent
		err = db.Where("id in (?)", studentIds).Find(&students).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.StudentMap = make(map[uint64]*types.ModelStudent, len(students))
		for _, s := range students {
			rsp.StudentMap[s.Id] = s
		}
	}

	return &rsp, nil
}

func ListPaperExaminer(ctx *goon.Ctx, req *types.ListPaperExaminerReq) (*types.ListPaperExaminerRsp, error) {
	var rsp types.ListPaperExaminerRsp

	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListPaper_OptionExaminerId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var err error
	rsp.PaperList, rsp.Page, err = paper.ListPaper(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if req.ShowExam && len(rsp.PaperList) > 0 {
		var exams []*types.ModelExam
		err = db.Where("id in (?)", utils.PluckUint64(rsp.PaperList, "ExamId")).Find(&exams).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.ExamMap = make(map[uint64]*types.ModelExam, len(exams))
		for _, exam := range exams {
			rsp.ExamMap[exam.Id] = exam
		}
	}

	if req.ShowStudent && len(rsp.PaperList) > 0 {
		var studentIds []uint64
		studentIds = append(studentIds, utils.PluckUint64(rsp.PaperList, "ExaminerId")...)
		studentIds = append(studentIds, utils.PluckUint64(rsp.PaperList, "ReviewerId")...)

		studentIds = pie.Unique(studentIds)

		var students []*types.ModelStudent
		err = db.Where("id in (?)", studentIds).Find(&students).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.StudentMap = make(map[uint64]*types.ModelStudent, len(students))
		for _, s := range students {
			rsp.StudentMap[s.Id] = s
		}
	}

	return &rsp, nil
}

func ListPaperReviewer(ctx *goon.Ctx, req *types.ListPaperExaminerReq) (*types.ListPaperExaminerRsp, error) {
	var rsp types.ListPaperExaminerRsp

	req.Options.Options = append(req.Options.Options, types.Option{
		Key: types.ListPaper_OptionReviewerId,
		Val: strconv.FormatUint(ctx.GetUint64(types.HeaderUserId), 10),
	})

	var err error
	rsp.PaperList, rsp.Page, err = paper.ListPaper(req.Options)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	if req.ShowExam && len(rsp.PaperList) > 0 {
		var exams []*types.ModelExam
		err = db.Where("id in (?)", utils.PluckUint64(rsp.PaperList, "ExamId")).Find(&exams).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.ExamMap = make(map[uint64]*types.ModelExam, len(exams))
		for _, exam := range exams {
			rsp.ExamMap[exam.Id] = exam
		}
	}

	if req.ShowStudent && len(rsp.PaperList) > 0 {
		var studentIds []uint64
		studentIds = append(studentIds, utils.PluckUint64(rsp.PaperList, "ExaminerId")...)
		studentIds = append(studentIds, utils.PluckUint64(rsp.PaperList, "ReviewerId")...)

		studentIds = pie.Unique(studentIds)

		var students []*types.ModelStudent
		err = db.Where("id in (?)", studentIds).Find(&students).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.StudentMap = make(map[uint64]*types.ModelStudent, len(students))
		for _, s := range students {
			rsp.StudentMap[s.Id] = s
		}
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
