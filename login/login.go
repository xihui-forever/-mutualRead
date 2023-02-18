package login

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/xtime"
	"github.com/xihui-forever/goon/middleware/session"
	"time"
)

var LoginHandlerMap = map[int]func(username interface{}, password string) (uint64, error){}

const (
	LoginTypeAdmin = iota + 1
	LoginTypeTeacher
	LoginTypeStudent
)

type (
	LoginReq struct {
		Logintype int         `json:"logintype,omitempty"`
		Username  interface{} `json:"username,omitempty"`
		Password  string      `json:"password,omitempty"`
	}
	LoginRsp struct {
		Token  string        `json:"token,omitempty"`
		Expire time.Duration `json:"expire,omitempty"`
	}
	LoginSession struct {
		role int    `json:"role,omitempty"`
		Id   uint64 `json:"id,omitempty"`
	}
)

func LoginHandler(req *LoginReq, session *session.Session) (*LoginRsp, error) {
	var resp LoginRsp
	id, err := LoginHandlerMap[req.Logintype](req.Username, req.Password)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Token, err = session.GenSession(&LoginSession{
		role: req.Logintype,
		Id:   id,
	}, xtime.Day)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	resp.Expire = xtime.Day
	return &resp, nil
}
