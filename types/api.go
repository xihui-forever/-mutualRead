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

const (
	ErrRolePermAlreadyExists ErrCode = iota + 10001
	ErrAdminExist
	ErrAdminNotExist
	ErrAdminRemoveFailed
	ErrPasswordWrong
	ErrPasswordChangeFailed
	ErrOneAppealAlreadyProgress
	ErrAppealAlreadyExists
	ErrAppealNotExists
	ErrAppealAlreadyHanded
	ErrExamNotExist
	ErrExamChangeFailed

	ErrPasswordIncorrect
	ErrTeacherNotExist

	ErrLoginTypeNotFound

	ErrPaperExist
	ErrPaperChangeFailed
	ErrPaperNotExist
	ErrGradeChangeFailed

	ErrStudentExist
	ErrStudentNotExist
	ErrStudentRemoveFailed
	ErrorEmailNoChange
	ErrEmailChangeFailed

	ErrTeacherExist
	ErrTeacherRemoveFailed
	ErrorNewPwdEmpty
	ErrEmailEmpty
	ErrPathAlreadyExists
)

var errMap = map[ErrCode]string{
	SysError:                 "系统错误",
	ErrInvalidParam:          "参数错误",
	ErrInvalidOptionKey:      "选项错误",
	ErrRolePermAlreadyExists: "rolePerm already exists",

	ErrAdminExist:           "管理员已存在",
	ErrAdminNotExist:        "管理员不存在",
	ErrAdminRemoveFailed:    "管理员删除失败",
	ErrPasswordWrong:        "密码错误",
	ErrPasswordChangeFailed: "密码修改失败",

	ErrOneAppealAlreadyProgress: "已有一个申诉正在进行中",
	ErrAppealAlreadyExists:      "申诉已存在",
	ErrAppealNotExists:          "申诉不存在",
	ErrAppealAlreadyHanded:      "该申诉无法撤销",

	ErrExamNotExist:     "考试不存在",
	ErrExamChangeFailed: "考试修改失败",

	ErrPasswordIncorrect: "密码错误",
	ErrTeacherNotExist:   "教师不存在",

	ErrLoginTypeNotFound: "登录类型不存在",

	ErrPaperExist:        "试卷已存在",
	ErrPaperChangeFailed: "试卷修改失败",
	ErrPaperNotExist:     "试卷不存在",
	ErrGradeChangeFailed: "成绩修改失败",

	ErrStudentExist:        "学生已存在",
	ErrStudentNotExist:     "学生不存在",
	ErrStudentRemoveFailed: "学生删除失败",
	ErrorEmailNoChange:     "邮箱没有变更",
	ErrEmailChangeFailed:   "邮箱修改失败",

	ErrTeacherExist:        "教师已存在",
	ErrTeacherRemoveFailed: "教师删除失败",
	ErrorNewPwdEmpty:       "新密码为空",
	ErrEmailEmpty:          "邮箱为空",

	ErrPathAlreadyExists: "路径已存在",
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
