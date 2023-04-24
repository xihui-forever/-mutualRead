package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ModelStudent struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_student_student_id,unique"`

	StudentId string `json:"student_id,omitempty" gorm:"column:student_id;not null;index:idx_student_student_id,unique"`
	Password  string `json:"password,omitempty" gorm:"column:password;not null"`
	Name      string `json:"name,omitempty" gorm:"column:name;not null"`
	Email     string `json:"email,omitempty" gorm:"column:email;not null"`
}

func (m *ModelStudent) GetId() uint64 {
	return m.Id
}

func (m *ModelStudent) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelStudent) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelStudent) TableName() string {
	return "goon_student"
}

const (
	CmdPathGetStudent    = "/student/get"
	CmdPathChangeStudent = "/student/change"

	CmdPathGetStudentAdmin  = "/admin/student/get"
	CmdPathAddStudentAdmin  = "/admin/student/set"
	CmdPathSetStudentAdmin  = "/admin/student/set"
	CmdPathDelStudentAdmin  = "/admin/student/del"
	CmdPathListStudentAdmin = "/admin/student/list"
)

type (
	GetStudentRsp struct {
		Student *ModelStudent `json:"student,omitempty" yaml:"student,omitempty"`
	}
)

const (
	StudentChangeTypeEmail = iota + 1
)

type (
	ChangeStudentReq struct {
		ChangeType int `json:"change_type,omitempty" yaml:"change_type,omitempty" validate:"required"`

		Email string `json:"email,omitempty" yaml:"email,omitempty"`
	}
)

type (
	AddStudentAdminReq struct {
		Student *ModelStudent `json:"student,omitempty" yaml:"student,omitempty" validate:"required"`
	}

	AddStudentAdminRsp struct {
		Student *ModelStudent `json:"student,omitempty" yaml:"student,omitempty"`
	}
)

type (
	GetStudentAdminReq struct {
		Id uint64 `json:"id,omitempty" yaml:"id,omitempty" validate:"required"`
	}

	GetStudentAdminRsp struct {
		Student *ModelStudent `json:"student,omitempty" yaml:"student,omitempty"`
	}
)

type (
	SetStudentAdminReq struct {
		Student *ModelStudent `json:"student,omitempty" yaml:"student,omitempty" validate:"required"`
	}

	SetStudentAdminRsp struct {
		Student *ModelStudent `json:"student,omitempty" yaml:"student,omitempty"`
	}
)

type (
	DelStudentAdminReq struct {
		Id uint64 `json:"id,omitempty" yaml:"id,omitempty"`
	}
)

const (
	ListStudent_OptionStudentId = iota + 1
	ListStudent_OptionNameLike
	ListStudent_OptionEmailLike
)

type (
	ListStudentAdminReq struct {
		Options *ListOption `json:"options,omitempty" yaml:"options,omitempty" validate:"required"`
	}

	ListStudentAdminRsp struct {
		Page     *Page           `json:"page,omitempty" yaml:"page,omitempty"`
		Students []*ModelStudent `json:"students,omitempty" yaml:"students,omitempty"`
	}
)
