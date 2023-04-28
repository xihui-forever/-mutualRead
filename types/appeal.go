package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

const (
	AppealStateWaitReviewer = iota + 1 // 等待阅卷人审核
	AppealStateWaitTeacher             // 等待老师审核
	AppealStateFinish                  // 申诉结束
	AppealStateRecall                  // 申诉撤回
)

type ModelAppeal struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null"`

	State   int    `json:"state,omitempty" gorm:"column:state;not null"`
	ExamId  uint64 `json:"exam_id,omitempty" gorm:"column:exam_id;not null"`
	PaperId uint64 `json:"paper_id,omitempty" gorm:"column:paper_id;not null"`

	ExaminerId uint64 `json:"examiner_id,omitempty" gorm:"column:examiner_id;not null"`
	ReviewerId uint64 `json:"reviewer_id,omitempty" gorm:"column:reviewer_id;not null"`
	TeacherId  uint64 `json:"teacher_id,omitempty" gorm:"column:teacher_id;not null"`

	ChangeAt uint32 `json:"change_at,omitempty" gorm:"column:change_at"`
	ReviewAt uint32 `json:"reviewer_at,omitempty" gorm:"column:reviewer_at"`
	ResultAt uint32 `json:"result_at,omitempty" gorm:"column:result_at"`

	AppealInfo   string `json:"appeal_info,omitempty" gorm:"column:appeal_info;not null"`
	ReviewInfo   string `json:"review_info,omitempty" gorm:"column:review_info"`
	AppealResult string `json:"appeal_result,omitempty" gorm:"column:appeal_result"`

	Grade int32 `json:"grade,omitempty" gorm:"column:grade;not null"`
}

func (m *ModelAppeal) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelAppeal) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelAppeal) TableName() string {
	return "goon_appeal"
}

const (
	CmdPathAddAppeal          = "/appeal/add"
	CmdPathListAppealTeacher  = "/teacher/appeal/list"
	CmdPathListAppealExaminer = "/examiner/appeal/list"
	CmdPathListAppealReviewer = "/reviewer/appeal/list"
	CmdPathGetAppealTeacher   = "/teacher/appeal/get"
	CmdPathGetAppealExaminer  = "/examiner/appeal/get"
	CmdPathGetAppealReviewer  = "/reviewer/appeal/get"

	CmdPathSetAppealExaminer = "/examiner/appeal/set"
	CmdPathSetAppealReviewer = "/reviewer/appeal/set"
	CmdPathSetAppealTeacher  = "/teacher/appeal/set"
	CmdPathRecallAppeal      = "/appeal/recall"
)

type (
	AddAppealReq struct {
		PaperId    uint64 `json:"paper_id,omitempty" validate:"required"`
		AppealInfo string `json:"appeal_info,omitempty" validate:"required"`
	}

	AddAppealRsp struct {
		Appeal *ModelAppeal `json:"appeal,omitempty"`
	}
)

const (
	ListAppeal_OptionTeacherId = iota + 1
	ListAppeal_OptionExaminerId
	ListAppeal_OptionReviewerId
	ListAppeal_OptionPaperId
	ListAppeal_OptionId
	ListAppeal_OptionStates
)

type (
	ListAppealTeacherReq struct {
		Options   *ListOption `json:"options,omitempty" validate:"required"`
		ShowPaper bool        `json:"show_paper"`
	}

	ListAppealTeacherRsp struct {
		Page     *Page                  `json:"page,omitempty"`
		List     []*ModelAppeal         `json:"list,omitempty"`
		PaperMap map[uint64]*ModelPaper `json:"paper_map,omitempty"`
	}
)

type (
	ListAppealExaminerReq struct {
		Options   *ListOption `json:"options,omitempty" validate:"required"`
		ShowPaper bool        `json:"show_paper"`
	}

	ListAppealExaminerRsp struct {
		Page     *Page                  `json:"page,omitempty"`
		List     []*ModelAppeal         `json:"list,omitempty"`
		PaperMap map[uint64]*ModelPaper `json:"paper_map,omitempty"`
	}
)

type (
	ListAppealReviewerReq struct {
		Options   *ListOption `json:"options,omitempty" validate:"required"`
		ShowPaper bool        `json:"show_paper"`
	}

	ListAppealReviewerRsp struct {
		Page     *Page                  `json:"page,omitempty"`
		List     []*ModelAppeal         `json:"list,omitempty"`
		PaperMap map[uint64]*ModelPaper `json:"paper_map,omitempty"`
	}
)

type (
	GetAppealTeacherReq struct {
		Id        uint64 `json:"id,omitempty" validate:"required"`
		ShowPaper bool   `json:"show_paper"`
	}

	GetAppealTeacherRsp struct {
		Appeal *ModelAppeal `json:"appeal,omitempty"`
		Paper  *ModelPaper  `json:"paper,omitempty"`
	}
)

type (
	GetAppealExaminerReq struct {
		Id        uint64 `json:"id,omitempty" validate:"required"`
		ShowPaper bool   `json:"show_paper"`
	}

	GetAppealExaminerRsp struct {
		Appeal *ModelAppeal `json:"appeal,omitempty"`
		Paper  *ModelPaper  `json:"paper,omitempty"`
	}
)

type (
	GetAppealReviewerReq struct {
		Id        uint64 `json:"id,omitempty" validate:"required"`
		ShowPaper bool   `json:"show_paper"`
	}

	GetAppealReviewerRsp struct {
		Appeal *ModelAppeal `json:"appeal,omitempty"`
		Paper  *ModelPaper  `json:"paper,omitempty"`
	}
)

type (
	SetAppealExaminerReq struct {
		AppealId   uint64 `json:"appeal_id,omitempty" validate:"required"`
		AppealInfo string `json:"appeal_info,omitempty"`
	}
)

type (
	SetAppealReviewerReq struct {
		AppealId   uint64 `json:"appeal_id,omitempty" validate:"required"`
		ReviewInfo string `json:"review_info,omitempty"`
	}
)

type (
	SetAppealTeacherReq struct {
		AppealId     uint64 `json:"appeal_id,omitempty" validate:"required"`
		Grade        int32  `json:"grade,omitempty" validate:"required"`
		AppealResult string `json:"appeal_result,omitempty"`
	}
)

type (
	RecallAppealReq struct {
		AppealId uint64 `json:"appeal_id,omitempty" validate:"required"`
	}
)

type (
	EventAppealChanged struct {
		Appeal *ModelAppeal `json:"appeal,omitempty" yaml:"appeal,omitempty"`
	}
)
