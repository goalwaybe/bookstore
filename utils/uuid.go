package utils

import "github.com/google/uuid"

// GenerateUUID 生成一个全局唯一的 UUID 字符串
func GenerateUUID() string {
	return uuid.New().String() // 返回形如 "3f29b8b0-2e6e-4a99-9f3b-5d0d3f8cfcf1"
}
