package models

import (
	"time"
)

type StepMdl struct {
	Id          int64     `gorm:"Column:id" json:"Id"`
	ParentId    int64     `gorm:"Column:parentid" json:"ParentId"`
	StepNumber  int       `gorm:"Column:stepnumber" json:"StepNumber"`
	Title       string    `gorm:"Column:title" json:"Title"`
	PlanStart   time.Time `gorm:"Column:planstart" json:"PlanStart"`
	PlanEnd     time.Time `gorm:"Column:planend" json:"PlanEnd"`
	Description string    `gorm:"Column:description" json:"Description"`
	Summary     string    `gorm:"Column:summary" json:"Summary"`
}

func (StepMdl) TableName() string {
	return "step"
}
