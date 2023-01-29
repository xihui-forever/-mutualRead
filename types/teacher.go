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
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_teacher_teacherId,unique"`

	TeacherId string `json:"teacherId,omitempty" gorm:"column:teacherId;not null;index:idx_teacher_teacherId,unique"`
	Password  string `json:"password,omitempty" gorm:"column:password;not null"`
	Name      string `json:"name,omitempty" gorm:"column:name;not null"`
	Email     string `json:"email,omitempty" gorm:"column:email;not null"`
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
