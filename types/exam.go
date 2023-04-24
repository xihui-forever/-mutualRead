package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ModelExam struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null"`

	Name      string `json:"name,omitempty" gorm:"column:name;not null"`
	TeacherId uint64 `json:"teacher_id,omitempty" gorm:"column:teacher_id;not null"`
}

func (m *ModelExam) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelExam) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelExam) TableName() string {
	return "goon_exam"
}

const (
	CmdPathAddExam  = "/exam/add"
	CmdPathSetExam  = "/exam/set"
	CmdPathDelExam  = "/exam/del"
	CmdPathListExam = "/exam/list"
)

type (
	AddExamReq struct {
		Exam *ModelExam `json:"exam,omitempty" yaml:"exam,omitempty" validate:"required"`
	}

	AddExamRsp struct {
		Exam *ModelExam `json:"exam,omitempty" yaml:"exam,omitempty" validate:"required"`
	}
)

type (
	SetExamReq struct {
		Exam *ModelExam `json:"exam,omitempty" yaml:"exam,omitempty" validate:"required"`
	}

	SetExamRsp struct {
		Exam *ModelExam `json:"exam,omitempty" yaml:"exam,omitempty" validate:"required"`
	}
)

type (
	DelExamReq struct {
		Id *ModelExam `json:"id,omitempty" yaml:"id,omitempty" validate:"required"`
	}
)

const (
	ListExam_OptionTeacherId = iota + 1
	ListExam_OptionNameLike
)

type (
	ListExamReq struct {
		Options *ListOption `json:"options,omitempty" yaml:"options,omitempty" validate:"required"`
	}

	ListExamRsp struct {
		Page  *Page        `json:"page,omitempty" yaml:"page,omitempty"`
		Exams []*ModelExam `json:"exams,omitempty" yaml:"exams,omitempty"`
	}
)
