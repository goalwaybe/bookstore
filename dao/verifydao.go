package dao

import "time"

// SaveResetCode 保存验证码到 Redis
func SaveResetCode(email, code string, expire time.Duration) error {
	return SetCache("reset:"+email, code, expire)
}

// GetResetCode 获取验证码
func GetResetCode(email string) (string, error) {
	return GetCache("reset:" + email)
}

// DeleteResetCode 删除验证码
func DeleteResetCode(email string) error {
	return DelCache("reset:" + email)
}
