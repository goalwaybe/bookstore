package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	cfg := Conf.Redis
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,         //Redis 地址
		Password: cfg.Password, //没有密码可以为空
		DB:       cfg.DB,       // 默认 DB
		PoolSize: cfg.PoolSize,
	})

	//测试连接
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败:%v", err)
	}
	log.Println("Redis 连接成功")
}
