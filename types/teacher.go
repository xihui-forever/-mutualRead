package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ModelTeacher struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_teacher_teacher_id,unique"`

	TeacherId string `json:"teacher_id,omitempty" gorm:"column:teacher_id;not null;index:idx_teacher_teacher_id,unique"`
	Password  string `json:"password,omitempty" gorm:"column:password;not null"`
	Name      string `json:"name,omitempty" gorm:"column:name;not null"`
	Email     string `json:"email,omitempty" gorm:"column:email;not null"`
}

func (m *ModelTeacher) GetId() uint64 {
	return m.Id
}

func (m *ModelTeacher) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelTeacher) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelTeacher) TableName() string {
	return "goon_teacher"
}

// cmd list
const (
	CmdPathGetTeacher    = "/teacher/get"
	CmdPathChangeTeacher = "/teacher/change"

	CmdPathGetTeacherAdmin  = "/admin/teacher/get"
	CmdPathAddTeacherAdmin  = "/admin/teacher/set"
	CmdPathSetTeacherAdmin  = "/admin/teacher/set"
	CmdPathDelTeacherAdmin  = "/admin/teacher/del"
	CmdPathListTeacherAdmin = "/admin/teacher/list"
)

const (
	TeacherChangeTypeEmail = iota + 1
)

type (
	GetTeacherRsp struct {
		Teacher *ModelTeacher `json:"teacher,omitempty" yaml:"teacher,omitempty"`
	}
)

type (
	ChangeTeacherReq struct {
		ChangeType int `json:"change_type,omitempty" validate:"required"`

		Email string `json:"email,omitempty"`
	}
)

type (
	GetTeacherAdminReq struct {
		Id uint64 `json:"id,omitempty" validate:"required"`
	}

	GetTeacherAdminRsp struct {
		Teacher *ModelTeacher `json:"teacher,omitempty" yaml:"teacher,omitempty"`
	}
)

type (
	AddTeacherAdminReq struct {
		Teacher *ModelTeacher `json:"teacher,omitempty" yaml:"teacher,omitempty" validate:"required"`
	}

	AddTeacherAdminRsp struct {
		Teacher *ModelTeacher `json:"teacher,omitempty" yaml:"teacher,omitempty"`
	}
)

type (
	SetTeacherAdminReq struct {
		Teacher *ModelTeacher `json:"teacher,omitempty" yaml:"teacher,omitempty" validate:"required"`
	}

	SetTeacherAdminRsp struct {
		Teacher *ModelTeacher `json:"teacher,omitempty" yaml:"teacher,omitempty"`
	}
)

type (
	DelTeacherAdminReq struct {
		Id uint64 `json:"id,omitempty"`
	}
)

const (
	ListTeacher_OptionTeacherId = iota + 1
	ListTeacher_OptionNameLike
	ListTeacher_OptionEmailLike
)

type (
	ListTeacherAdminReq struct {
		Options *ListOption `json:"options,omitempty" validate:"required"`
	}

	ListTeacherAdminRsp struct {
		Page     *Page           `json:"page,omitempty"`
		Teachers []*ModelTeacher `json:"teachers,omitempty" yaml:"teachers,omitempty"`
	}
)
