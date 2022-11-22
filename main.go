package main

import (
	"GoIM/router"
	"GoIM/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()

	r := router.Router()

	r.Run() // listen and serve on 0.0.0.0:8080
}
