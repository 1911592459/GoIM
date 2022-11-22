package main

import (
	"GoIM/models"
	"GoIM/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	// 迁移 schema
	//utils.DB.AutoMigrate(&models.UserBasic{})
	utils.DB.AutoMigrate(&models.Message{})
	utils.DB.AutoMigrate(&models.Contact{})
	utils.DB.AutoMigrate(&models.GroupBasic{})
	/*	// Create
		user := &models.UserBasic{}
		user.Name = "申专"
		utils.DB.Create(user)

		// Read
		fmt.Println(utils.DB.First(user, 6)) // 根据整型主键查找

		// Update - 将 product 的 price 更新为 200
		utils.DB.Model(user).Update("PassWord", "1234")*/

}
