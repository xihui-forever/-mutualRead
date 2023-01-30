package types

import (
	"gorm.io/plugin/soft_delete"
)

type ModelAppeal struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null"`

	PaperId      uint64 `json:"paperId,omitempty" gorm:"column:paperId;not null"`
	AppealStatus string `json:"appealStatus,omitempty" gorm:"column:appealStatus;not null"`
	AppealInfo   string `json:"password,omitempty" gorm:"column:password"`
	ReviewInfo   string `json:"password,omitempty" gorm:"column:reviewInfo"`
	AppealResult string `json:"appealResult,omitempty" gorm:"column:appealResult"`
}
