package models

type FlagMdl struct {
	Id          int64  `gorm:"Column:id" json:"Id"`
	UserId      int64  `gorm:"Column:userid" json:"UserId"`
	Name        string `gorm:"Column:name" json:"Name"`
	Description string `gorm:"Column:description" json:"Description"`
	Summary     string `gorm:"Column:summary" json:"Summary"`
	CreateAt    string `gorm:"Column:createat" json:"CreateAt"`
	UpdateAt    string `gorm:"Column:updateat" json:"UpdateAt"`
}

func (FlagMdl) TableName() string {
	return "flag"
}
