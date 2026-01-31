package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

/**
Redis 测试（_test.go 格式)
*/

func TestRedisConnection(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("redis 连接失败： %v", err)
	}
	t.Logf("Redis 连接成功:%v", pong)

	//测试写入
	err = rdb.Set(ctx, "test_key", "Hello Redis", 10*time.Minute).Err()
	if err != nil {
		t.Fatalf("Redis 写入失败: %v", err)
	}

	val, err := rdb.Get(ctx, "test_key").Result()
	if err != nil {
		t.Fatalf("Redis 读取失败：%v", err)
	}
	if val != "Hello Redis" {
		t.Fatalf("Redis 值不正确，got:%v", val)
	}
	t.Logf("Redis 读取值： %v", val)
}
