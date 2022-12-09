package main

import (
	"GoIM/models"
	"GoIM/router"
	"GoIM/utils"
	"github.com/spf13/viper"
	"time"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	InitTimer()
	r := router.Router()
	r.Run(viper.GetString("port.server"))
}

//初始化定时器
func InitTimer() {
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second, time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second, models.CleanConnection, "")
}
