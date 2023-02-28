package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ModelPerm struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_perm_role,unique"`

	Role       int    `json:"role,omitempty" gorm:"column:role;not null;index:idx_perm_role,unique"`
	Permission string `json:"permission,omitempty" gorm:"column:permission;not null"`
}

func (m *ModelPerm) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelPerm) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelPerm) TableName() string {
	return "mutual_read_perm"
}
