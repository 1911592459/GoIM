package models

import "gorm.io/gorm"

//群信息
type GroupBasic struct {
	gorm.Model
	Name    string //群名称
	OwnerId uint   //群主
	Icon    string //图片
	Desc    string //描述
	Type    int
}

func (table *GroupBasic) TableName() string {
	return "group_Basic"
}
