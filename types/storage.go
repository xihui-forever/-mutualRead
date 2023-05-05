package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"gorm.io/plugin/soft_delete"
)

type ResourceDetail struct {
	Body []byte `json:"body,omitempty"`

	ContentType string `json:"content_type"`
}

func (m *ResourceDetail) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ResourceDetail) Value() (driver.Value, error) {
	return utils.Value(m)
}

type ModelResource struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null;index:idx_resource_path,unique"`

	Path string `json:"path,omitempty" gorm:"column:path;not null;index:idx_resource_path,unique"`

	ResourceDetail *ResourceDetail `json:"resource_detail" gorm:"column:resource_detail;type:json"`
}

func (*ModelResource) TableName() string {
	return "goon_resource"
}

const (
	CmdPathResourcePut = "/resource/put"
	CmdPathResourceGet = "/resource/get/"
)

type PutResourceReq struct {
	Body []byte `json:"body,omitempty"`

	ContentType string `json:"content_type"`

	Path string `json:"path,omitempty"`
}

type PutResourceRsp struct {
	Path string `json:"path,omitempty"`
}
