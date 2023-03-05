package types

type ErrCode int32

const (
	SysError ErrCode = -1

	Success ErrCode = 0
)

type Error struct {
	Code ErrCode `json:"code" yaml:"code,omitempty"`
	Msg  string  `json:"msg" yaml:"msg,omitempty"`
}

func (p *Error) Error() string {
	if p == nil {
		return ""
	}

	return p.Msg
}
