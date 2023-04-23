package impl

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/mutualRead/exam"
	"github.com/xihui-forever/mutualRead/types"
)

func init() {
	CmdList = append(CmdList, Cmd{
		Path:  "/exam_list",
		Role:  2,
		Logic: GetExamList,
	}, Cmd{
		Path:  "/exam_get",
		Role:  2,
		Logic: GetTeacher,
	}, Cmd{
		Path:  "/exam_add",
		Role:  2,
		Logic: ChangeTeacher,
	}, Cmd{
		Path:  "/exam_del",
		Role:  2,
		Logic: GetExamList,
	})
}

type (
	ExamListReq struct {
		TeacherId string `json:"teacher_id"`
	}
	ExamListRsp struct {
		ExamList *[]types.ModelExam
	}
	GetExamReq struct {
		Id uint64 `json:"id"`
	}
	AddExamReq struct {
		Name      string `json:"name"`
		TeacherId string `json:"teacher_id"`
	}
	DelExamReq struct {
		Id        uint64 `json:"id"`
		TeacherId string `json:"teacher_id"`
	}
)

func GetExamList(ctx *goon.Ctx, req *ExamListReq) (*ExamListRsp, error) {
	var res ExamListRsp
	a, err := exam.GetExamList(req.TeacherId)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	res.ExamList = a

	return &res, nil
}

func GetExam(ctx *goon.Ctx, req *GetExamReq) (*types.ModelExam, error) {
	a, err := exam.GetExam(req.Id)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return a, nil
}

func AddExam(ctx *goon.Ctx, req *AddExamReq) (*ExamListRsp, error) {
	var a *ExamListRsp
	_, err := exam.AddExam(req.Name, req.TeacherId)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	a.ExamList, err = exam.GetExamList(req.TeacherId)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return a, nil
}

func DelExam(ctx *goon.Ctx, req *DelExamReq) (*ExamListRsp, error) {
	var a *ExamListRsp
	err := exam.RemoveExam(req.Id)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	a.ExamList, err = exam.GetExamList(req.TeacherId)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return a, nil
}
