package utils

import "net/http"

// GetUserID 安全获取用户ID
func GetUserID(r *http.Request) int {
	uid := r.Context().Value("userID")
	if uid == nil {
		return 0
	}
	id, ok := uid.(int)
	if !ok {
		return 0
	}
	return id
}
