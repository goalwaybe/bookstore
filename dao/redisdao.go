package dao

import (
	"bookstore/config"
	"bookstore/utils"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetCache 设置缓存
func SetCache(key string, value interface{}, expiration time.Duration) error {
	return config.RedisClient.Set(config.Ctx, key, value, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string) (string, error) {
	return config.RedisClient.Get(config.Ctx, key).Result()
}

// DelCache 删除缓存
func DelCache(key string) error {
	return config.RedisClient.Del(config.Ctx, key).Err()
}

// BlacklistToken 添加 token 到黑名单
func BlacklistToken(tokenID string, expire time.Duration) error {
	return SetCache("blacklist:"+tokenID, "blacklisted", expire)
}

// IsTokenBlacklisted 检查 token 是否在黑名单
//func IsTokenBlacklisted(jti string) (bool, error) {
//	val, err := GetCache("blacklist:" + jti)
//	if err != nil && err.Error() != "redis:nil" {
//		log.Printf("检查黑名单时出错：%v", err)
//		return false, err
//	}
//	return val != "", nil
//}

func IsTokenBlacklisted(jti string) (bool, error) {
	val, err := GetCache("blacklist:" + jti)
	if errors.Is(err, redis.Nil) {
		// ⚙️ key 不存在，说明没被拉黑，正常
		return false, nil
	}
	if err != nil {
		log.Printf("检查黑名单时出错：%v", err)
		return false, err
	}
	return val == "blacklisted", nil
}

//var ctx = context.Background()

// TryLock 尝试获取一个分布式锁
// 返回 true 表示加锁成功: false 表示锁已存在
func TryLock(key string, expiration time.Duration) (bool, error) {
	ok, err := config.RedisClient.SetNX(ctx, key, "locked", expiration).Result()
	return ok, err
}

// Unlock 释放锁
func Unlock(key string) error {
	return config.RedisClient.Del(ctx, key).Err()
}

// RevokeTokenFromRequest 从请求头提取 Token 并加入黑名单
func RevokeTokenFromRequest(r *http.Request) error {
	// 从请求头获取用户ID
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return errors.New("未提供 Token")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return errors.New("authorization 格式错误")
	}

	tokenString := parts[1]

	claims, err := utils.ParseToken(tokenString)

	if err != nil {
		return errors.New("token 无效")
	}

	expireDuration := time.Until(claims.ExpiresAt.Time)

	if err := BlacklistToken(claims.ID, expireDuration); err != nil {
		return errors.New("token 注销失败")
	}
	return nil
}
