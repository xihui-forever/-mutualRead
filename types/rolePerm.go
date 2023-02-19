package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ModelRolePerm struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_rolePerm_role,unique"`

	Role       string `json:"role,omitempty" gorm:"column:role;not null;index:idx_rolePerm_role,unique"`
	Permission string `json:"permission,omitempty" gorm:"column:permission;not null"`
}

func (m *ModelRolePerm) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelRolePerm) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelRolePerm) TableName() string {
	return "goon_roleMap"
}
