package models

import (
	"time"
)

type UserMdl struct {
	Id         int64     `gorm:"Column:id" json:"Id"`
	Username   string    `gorm:"Column:username" json:"Username"`
	Icon       string    `gorm:"Column:icon" json:"Icon" json:"Icon"`
	Name       string    `gorm:"Column:name" json:"Name"`
	Password   []byte    `gorm:"Column:pwd" json:"-"`
	Mobile     string    `gorm:"Column:mobile" json:"Mobile"`
	Email      string    `gorm:"Column:email" json:"Email"`
	UUID       string    `gorm:"Column:uuid" json:"UUID"`
	State      int       `gorm:"Column:state" json:"State"`
	LockTime   time.Time `gorm:"Column:locktime" json:"LockTime"`
	LockReason string    `gorm:"Column:lockreason" json:"LockReason"`
	CreateAt   time.Time `gorm:"Column:createat" json:"CreateAt"`
	UpdateAt   time.Time `gorm:"Column:updateat" json:"UpdateAt"`
}

func (UserMdl) TableName() string {
	return "user"
}
