package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ModelPaper struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_paper_exam_id,unique"`

	ExamId     uint64 `json:"exam_id,omitempty" gorm:"column:exam_id;not null;index:idx_paper_exam_id,unique"`
	ExaminerId uint64 `json:"examiner_id,omitempty" gorm:"column:examiner_id;not null;index:idx_paper_exam_id,unique"`
	ReviewerId uint64 `json:"reviewer_id,omitempty" gorm:"column:reviewer_id;not null"`
	TeacherId  uint64 `json:"teacher_id,omitempty" gorm:"column:teacher_id;not null"`

	// 分数
	Grade  uint32 `json:"grade,omitempty" gorm:"column:grade;not null"`
	ImgUrl string `json:"img_url,omitempty" gorm:"column:img_url;not null"`
}

func (m *ModelPaper) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelPaper) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelPaper) TableName() string {
	return "goon_paper"
}

const (
	CmdPathAddPaper = "/paper/add"

	CmdPathListPaperTeacher  = "/teacher/paper/list"
	CmdPathListPaperExaminer = "/examiner/paper/list"
	CmdPathListPaperReviewer = "/reviewer/paper/list"

	CmdPathGetPaperTeacher  = "/teacher/paper/get"
	CmdPathGetPaperExaminer = "/examiner/paper/get"
	CmdPathGetPaperReviewer = "/reviewer/paper/get"
	CmdPathDelPaperTeacher  = "/paper/del"
)

type (
	AddPaperReq struct {
		ExamId     uint64 `json:"exam_id,omitempty" gorm:"column:exam_id;not null;index:idx_paper_exam_id,unique"`
		ExaminerId string `json:"examiner_id,omitempty" gorm:"column:examiner_id;not null;index:idx_paper_exam_id,unique"`
		ReviewerId string `json:"reviewer_id,omitempty" gorm:"column:reviewer_id;not null"`

		// 分数
		Grade  uint32 `json:"grade,omitempty" gorm:"column:grade;not null"`
		ImgUrl string `json:"img_url,omitempty" gorm:"column:img_url;not null"`
	}

	AddPaperRsp struct {
		Paper *ModelPaper `json:"paper,omitempty" validate:"required"`
	}
)

const (
	ListPaper_OptionExamId = iota + 1
	ListPaper_OptionExaminerId
	ListPaper_OptionReviewerId
	ListPaper_OptionTeacherId
	ListPaper_OptionId
)

type (
	ListPaperTeacherReq struct {
		Options     *ListOption `json:"options,omitempty"`
		ShowExam    bool        `json:"show_exam,omitempty"`
		ShowStudent bool        `json:"show_student,omitempty"`
	}

	ListPaperTeacherRsp struct {
		Page       *Page                    `json:"page,omitempty"`
		PaperList  []*ModelPaper            `json:"paper_list,omitempty"`
		ExamMap    map[uint64]*ModelExam    `json:"exam_map,omitempty"`
		StudentMap map[uint64]*ModelStudent `json:"student_map,omitempty"`
	}
)

type (
	ListPaperExaminerReq struct {
		Options     *ListOption `json:"options,omitempty"`
		ShowExam    bool        `json:"show_exam,omitempty"`
		ShowStudent bool        `json:"show_student,omitempty"`
	}

	ListPaperExaminerRsp struct {
		Page       *Page                    `json:"page,omitempty"`
		PaperList  []*ModelPaper            `json:"paper_list,omitempty"`
		ExamMap    map[uint64]*ModelExam    `json:"exam_map,omitempty"`
		StudentMap map[uint64]*ModelStudent `json:"student_map,omitempty"`
	}
)

type (
	GetPaperTeacherReq struct {
		Id uint64 `json:"id,omitempty" validate:"required"`
	}

	GetPaperTeacherRsp struct {
		Paper *ModelPaper `json:"paper,omitempty"`
	}
)
type (
	GetPaperExaminerReq struct {
		Id uint64 `json:"id,omitempty" validate:"required"`
	}

	GetPaperExaminerRsp struct {
		Paper *ModelPaper `json:"paper,omitempty"`
	}
)
type (
	GetPaperReviewerReq struct {
		Id uint64 `json:"id,omitempty" validate:"required"`
	}

	GetPaperReviewerRsp struct {
		Paper *ModelPaper `json:"paper,omitempty"`
	}
)

type (
	DelPaperTeacherReq struct {
		Id uint64 `json:"id,omitempty" validate:"required"`
	}
)
