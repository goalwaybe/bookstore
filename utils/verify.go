package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateVerifyCode 生成指定长度的数字验证码(安全随机)
func GenerateVerifyCode(length int) string {
	const digits = "0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		// 使用 crypto/rand 确保安全随机
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		code[i] = digits[n.Int64()]
	}
	return string(code)
}
