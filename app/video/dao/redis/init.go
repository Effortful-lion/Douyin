package redis

import (
	"Douyin/config"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	// 初始化 Redis 客户端
	host := config.Config.RedisConfig.RedisHost
	port := config.Config.RedisConfig.RedisPort
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),        // Redis 服务器地址
		Password: config.Config.RedisConfig.RedisPassword, // Redis 访问密码（如果有的话）
		DB:       1,                                       // Redis 数据库索引
	})
	// 测试连接
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("redis初始化成功...")
}
