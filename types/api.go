package types

type ErrCode int32

const (
	SysError ErrCode = -1

	Success ErrCode = 0
)

const (
	ErrInvalidParam ErrCode = iota + 1001
	ErrInvalidOptionKey
)

var errMap = map[ErrCode]string{
	SysError:            "系统错误",
	ErrInvalidParam:     "参数错误",
	ErrInvalidOptionKey: "选项错误",
}

type Error struct {
	Code ErrCode `json:"code" yaml:"code,omitempty"`
	Msg  string  `json:"msg" yaml:"msg,omitempty"`

	Data     any    `json:"data,omitempty" yaml:"data,omitempty"`
	TranceId string `json:"trance_id,omitempty" yaml:"trance_id,omitempty"`
}

func (p *Error) Error() string {
	if p == nil {
		return ""
	}

	return p.Msg
}

func CreateError(code ErrCode) error {
	if msg, ok := errMap[code]; ok {
		return CreateErrorWithMsg(code, msg)
	} else {
		return CreateErrorWithMsg(code, "")
	}
}

func CreateErrorWithMsg(code ErrCode, msg string) error {
	if code == Success {
		return nil
	}

	return &Error{
		Code: code,
		Msg:  msg,
	}
}
