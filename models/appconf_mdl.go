package models

type AppConfMdl struct {
	Id    int64  `gorm:"Column:id" json:"Id"`
	Name  string `gorm:"Column:name" json:"Name"`
	Value string `gorm:"Column:value" json:"Value"`
	State int    `gorm:"Column:state" json:"State"`
}

func (AppConfMdl) TableName() string {
	return "appconf"
}
