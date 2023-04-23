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
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null"`

	Img        string `json:"img,omitempty" gorm:"column:img;not null"`
	Grade      uint32 `json:"grade,omitempty" gorm:"column:grade;not null"`
	ExamId     uint64 `json:"exam_id,omitempty" gorm:"column:exam_id;not null;index:idx_paper_exam_id,unique"`
	ExaminerId uint64 `json:"examiner_id,omitempty" gorm:"column:examiner_id;not null;index:idx_paper_exam_id,unique"`
	ReviewerId uint64 `json:"reviewer_id,omitempty" gorm:"column:reviewer_id;not null"`
	TeacherId  uint64 `json:"teacher_id,omitempty" gorm:"column:teacher_id;not null"`
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
