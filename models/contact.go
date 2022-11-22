package models

import (
	"GoIM/utils"
	"gorm.io/gorm"
)

//人员关系，owner和tagrget的关系信息
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint //对于的谁
	Type     int  //类型， 1好友, 2群友，3

}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}
