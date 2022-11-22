package models

import (
	"GoIM/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name     string
	PassWord string
	//手机号，使用正则匹配
	Phone string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	//邮箱格式使用第三方工具类
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time
	IsLoginOut    bool
	DeviceInfo    string
}

func FindUserByNameAndPwd(name, pwd string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("name =? and pass_word =?", name, pwd).First(user)
	//token 加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.Md5Encode(str)
	utils.DB.Model(user).Where("id =?", user.ID).Update("identity", temp)
	return user
}
func FindUserByName(name string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("name =?", name).First(user)
	return user
}
func FindUserByPhone(phone string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("Phone =?", phone).First(user)
	return user
}
func FindUserByEmail(email string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("email =?", email).First(user)
	return user
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, datum := range data {
		fmt.Println(datum)
	}
	return data
}

func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}
func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(user)
}
func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
