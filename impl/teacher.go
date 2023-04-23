package impl

import (
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/teacher"
	"github.com/xihui-forever/mutualRead/types"
)

func init() {
	CmdList = append(CmdList, Cmd{
		Path:  "/teacher_get",
		Role:  2,
		Logic: GetTeacher,
	}, Cmd{
		Path:  "/teacher_change",
		Role:  2,
		Logic: ChangeTeacher,
	}, Cmd{
		Path:  "/teacher_exam_list",
		Role:  0,
		Logic: GetExamList,
	})
}

type (
	GetReq struct {
		Token string `json:"token"`
	}
	ChangeReq struct {
		TeacherId string `json:"teacher_id"`
		Email     string `json:"email"`
	}
)

func GetTeacher(ctx *goon.Ctx, req *GetReq) (*types.ModelTeacher, error) {
	var a *types.ModelTeacher
	s, err := session.GetSession(req.Token)
	if err != nil {
		log.Errorf("err:%v", err)
		return a, err
	}

	var loginReq LoginSession
	err = sonic.UnmarshalString(s, &loginReq)
	if err != nil {
		log.Errorf("err:%v", err)
		return a, err
	}

	a, err = teacher.GetTeacherById(loginReq.Id)
	if err != nil {

		log.Errorf("err:%v", err)
		return nil, err
	}
	return a, nil
}

func ChangeTeacher(ctx *goon.Ctx, req *ChangeReq) (*types.ModelTeacher, error) {
	var a *types.ModelTeacher
	err := teacher.ChangeEmail(req.TeacherId, req.Email)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	a, err = teacher.GetTeacher(req.TeacherId)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return a, nil
}
