package models

import (
	"time"
)

type EventMdl struct {
	Id          int64     `gorm:"Column:id" json:"Id"`
	UserId      int64     `gorm:"Column:userid" json:"UserId"`
	HappenTime  time.Time `gorm:"Column:happentime" json:"HappenTime"`
	Place       string    `gorm:"Column:place" json:"Place"`
	ContentType string    `gorm:"Column:contenttype" json:"ContentType"`
	Content     string    `gorm:"Column:content" json:"Content"`
	Summary     string    `gorm:"Column:summary" json:"Summary"`
	CreateAt    time.Time `gorm:"Column:createat" json:"CreateAt"`
	UpdateAt    time.Time `gorm:"Column:updateat" json:"UpdateAt"`
}

func (EventMdl) TableName() string {
	return "event"
}
