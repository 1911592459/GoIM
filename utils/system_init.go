package utils

import (
	"context"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis2.Client
)

// InitConfig 使用viper第三方工具类读取配置文件信息
func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql:", viper.Get("app.mysql"))
}
func InitMysql() {
	//自定义日志模板，打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		})
	DB, _ = gorm.Open(mysql.Open(viper.GetString("app.mysql.dns")), &gorm.Config{Logger: newLogger})
	/*if err != nil {
		panic("failed to connect database")
	}*/

}

func InitRedis() {
	Red = redis2.NewClient(&redis2.Options{
		Addr: viper.GetString("app.redis.addr"),
		//Password:     viper.GetString("app.redis.password"),
		DB:           viper.GetInt("app.redis.DB"),
		PoolSize:     viper.GetInt("app.redis.poolSize"),
		MinIdleConns: viper.GetInt("app.redis.minIdleConn"),
	})
	/*result, err := Red.Ping().Result()
	if err != nil {
		fmt.Println("redis init failed :,", err)
	} else {
		fmt.Println("redis init success", result)
	}*/

}

const (
	PublishKey = "websocket"
)

//利用redis的发布订阅功能实现消息的订阅和发布

//发布消息到指定频道
func Publish(ctx context.Context, channel string, msg string) error {
	err := Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

//使用redis订阅频道,然后获取频道里的消息
func Subscribe(ctx context.Context, channel string) (msg string, err error) {
	sub := Red.Subscribe(ctx, channel)
	message, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return message.Payload, err
}
