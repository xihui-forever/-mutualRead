package types

import (
	"gorm.io/plugin/soft_delete"
)

type ModelAppeal struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null"`

	State        int    `json:"state,omitempty" gorm:"column:appeal_status;not null"`
	PaperId      uint64 `json:"paper_id,omitempty" gorm:"column:paper_id;not null"`
	ExaminerId   uint64 `json:"examiner_id,omitempty" gorm:"column:examiner_id;not null"`
	ReviewerId   uint64 `json:"reviewer_id,omitempty" gorm:"column:reviewer_id;not null"`
	TeacherId    uint64 `json:"teacher_id,omitempty" gorm:"column:teacher_id;not null"`
	ChangeAt     uint32 `json:"change_at,omitempty" gorm:"column:change_at"`
	ReviewAt     uint32 `json:"reviewer_at,omitempty" gorm:"column:reviewer_at"`
	ResultAt     uint32 `json:"result_at,omitempty" gorm:"column:result_at"`
	AppealInfo   string `json:"appeal_info,omitempty" gorm:"column:appeal_info;not null"`
	ReviewInfo   string `json:"review_info,omitempty" gorm:"column:review_info"`
	AppealResult string `json:"appeal_result,omitempty" gorm:"column:appeal_result"`
}
