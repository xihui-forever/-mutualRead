package impl

import (
	"errors"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/xtime"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/admin"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/student"
	"github.com/xihui-forever/mutualRead/teacher"
	"time"
)

func init() {
	LoginHandlerMap[LoginTypeTeacher] = func(username string, password string) (uint64, error) {
		data, err := teacher.GetTeacher(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		err = teacher.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		return data.Id, nil
	}
	LoginHandlerMap[LoginTypeAdmin] = func(username string, password string) (uint64, error) {
		data, err := admin.GetAdmin(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		err = admin.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		return data.Id, nil
	}
	LoginHandlerMap[LoginTypeStudent] = func(username string, password string) (uint64, error) {
		data, err := student.GetStudent(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		err = student.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return 0, err
		}

		return data.Id, nil
	}

	CmdList = append(CmdList, Cmd{
		Path:  "/login",
		Role:  role.RoleTypePublic,
		Logic: LoginHandler,
	})
}

var LoginHandlerMap = map[int]func(username string, password string) (uint64, error){}

const (
	LoginTypeAdmin = iota + 1
	LoginTypeTeacher
	LoginTypeStudent
)

type (
	LoginReq struct {
		RoleType int    `json:"role_type,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
	LoginRsp struct {
		Token  string        `json:"token,omitempty"`
		Expire time.Duration `json:"expire,omitempty"`
	}
	LoginSession struct {
		RoleType int    `json:"role_type,omitempty"`
		Id       uint64 `json:"id,omitempty"`
	}
)

func LoginHandler(ctx *goon.Ctx, req *LoginReq) (*LoginRsp, error) {
	var resp LoginRsp
	logic, ok := LoginHandlerMap[req.RoleType]
	if !ok {
		log.Errorf("login type %v not found", req.RoleType)
		return nil, errors.New("login type not found")
	}

	id, err := logic(req.Username, req.Password)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Token, err = session.GenSession(&LoginSession{
		RoleType: req.RoleType,
		Id:       id,
	}, xtime.Day)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Expire = xtime.Day
	return &resp, nil
}
