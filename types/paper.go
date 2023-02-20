package types

import (
	"database/sql/driver"
	"github.com/darabuchi/utils"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/image"
	"gorm.io/plugin/soft_delete"
)

type ModelPaper struct {
	Id uint64 `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true;column:id;not null"`

	CreatedAt uint32                `json:"created_at,omitempty" gorm:"autoCreateTime;<-:create;column:created_at;not null"`
	UpdatedAt uint32                `json:"updated_at,omitempty" gorm:"autoUpdateTime;<-;column:updated_at;not null"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at;not null"`

	Imag     image.Derived `json:"imag,omitempty" gorm:"column:imag;not null"`
	Grade    uint32        `json:"grade,omitempty" gorm:"column:grade;not null"`
	ExamId   uint64        `json:"exam_id,omitempty" gorm:"column:exam_id;not null;index:idx_paper_exam_id,unique"`
	Examiner uint64        `json:"examiner,omitempty" gorm:"column:examiner;not null;index:idx_paper_exam_id,unique"`
	Reviewer uint64        `json:"reviewer,omitempty" gorm:"column:reviewer;not null"`
}

func (m *ModelPaper) Scan(value interface{}) error {
	return utils.Scan(value, m)
}

func (m *ModelPaper) Value() (driver.Value, error) {
	return utils.Value(m)
}

func (m *ModelPaper) TableName() string {
	return "mutual_read_paper"
}
