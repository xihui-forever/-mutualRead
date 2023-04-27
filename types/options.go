package types

import (
	"github.com/darabuchi/log"
	"gorm.io/gorm"
)

type (
	Option struct {
		Key int `json:"key"`
		Val any `json:"val"`
	}

	ListOption struct {
		Options []Option `json:"options"`
		Offset  uint32   `json:"offset,omitempty"`
		Limit   uint32   `json:"limit,omitempty"`

		ShowTotal bool `json:"show_total,omitempty"`
	}

	Page struct {
		Offset uint32 `json:"offset,omitempty"`
		Limit  uint32 `json:"limit,omitempty"`

		Total int64 `json:"total,omitempty"`
	}
)

func (p *ListOption) Find(db *gorm.DB, list interface{}) (*Page, error) {
	if p.Limit == 0 {
		p.Limit = 100
	}

	page := &Page{
		Offset: p.Offset,
		Limit:  p.Limit,
		Total:  0,
	}

	if p.ShowTotal {
		err := db.Count(&page.Total).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}

	err := db.Offset(int(p.Offset)).Limit(int(p.Limit)).Find(list).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return page, nil
}
