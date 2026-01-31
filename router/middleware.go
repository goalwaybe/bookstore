package router

import (
	"bookstore/dao"
	"bookstore/utils"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// JSONResponse 是通用响应格式
type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// JWTAuthMiddleware 验证 JWT Token 的中间件
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 从 Authorization 头中获取 token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, "缺少 Authorization Header", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSON(w, http.StatusUnauthorized, "Authorization 格式错误", nil)
			return
		}

		tokenString := parts[1]

		// 解析并验证 token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, "Token 无效或已失效", nil)
			return
		}

		isBlacklisted, err := dao.IsTokenBlacklisted(claims.ID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, "Redis 查询失败", nil)
			return
		}

		if isBlacklisted {
			writeJSON(w, http.StatusUnauthorized, "Token 已失效", nil)
			return
		}

		//// 将用户信息放入上下文（如果后续需要）
		//r.Header.Set("X-User-ID", strconv.Itoa(claims.UserID))
		//r.Header.Set("x-User-Type", claims.UserType)

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userType", claims.UserType)

		// 继续执行后续处理
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// writeJSON 是通用的 JSON 输出函数
func writeJSON(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(JSONResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
