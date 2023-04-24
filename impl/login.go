package impl

import (
	"errors"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/xtime"
	"github.com/xihui-forever/goon"
	"github.com/xihui-forever/goon/middleware/session"
	"github.com/xihui-forever/mutualRead/admin"
	"github.com/xihui-forever/mutualRead/role"
	"github.com/xihui-forever/mutualRead/rpc"
	"github.com/xihui-forever/mutualRead/student"
	"github.com/xihui-forever/mutualRead/teacher"
	"github.com/xihui-forever/mutualRead/types"
	"time"
)

func init() {
	LoginHandlerMap[LoginTypeTeacher] = func(username string, password string) (LoginBaseRsp, error) {
		data, err := teacher.GetTeacher(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = teacher.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return data, nil
	}
	LoginHandlerMap[LoginTypeAdmin] = func(username string, password string) (LoginBaseRsp, error) {
		data, err := admin.GetAdmin(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = admin.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return data, nil
	}
	LoginHandlerMap[LoginTypeStudent] = func(username string, password string) (LoginBaseRsp, error) {
		data, err := student.GetStudent(username)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = student.CheckPassword(password, data.Password)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return data, nil
	}

	rpc.Register("/login", LoginHandler, role.RoleTypePublic)
}

type LoginBaseRsp interface {
	GetId() uint64
}

var LoginHandlerMap = map[int]func(username string, password string) (LoginBaseRsp, error){}

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
		Info   interface{}   `json:"info,omitempty"`
	}
)

func LoginHandler(ctx *goon.Ctx, req *LoginReq) (*LoginRsp, error) {
	var resp LoginRsp
	logic, ok := LoginHandlerMap[req.RoleType]
	if !ok {
		log.Errorf("login type %v not found", req.RoleType)
		return nil, errors.New("login type not found")
	}

	info, err := logic(req.Username, req.Password)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Token, err = session.GenSession(&types.LoginSession{
		RoleType: req.RoleType,
		Id:       info.GetId(),
	}, xtime.Day)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Expire = xtime.Day
	resp.Info = info
	return &resp, nil
}
