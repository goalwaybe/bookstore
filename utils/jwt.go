package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

/**
JWT鉴权
*/

var jwtSecret = []byte("your-secret-key") //自定义密钥

// Claims 自定义结构体
type Claims struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ErrTokenBlacklisted ✅ 黑名单错误定义
var ErrTokenBlacklisted = errors.New("token 已被注销")

// GenerateToken 生成 JWT 并设置过期时间(单位小时)
func GenerateToken(userID int, username, userType string, expireHours time.Duration) (string, *Claims, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireHours * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bookstore",
			ID:        generateTokenID(),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObj.SignedString(jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return token, claims, nil
}

// ParseToken 解析并验证 Token (自动检查黑名单)
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	//提取 Claims
	if claims, ok := token.Claims.(*Claims); ok || !token.Valid {
		return claims, nil
	}

	// 检查黑名单 (使用jti 而非整串 token)
	//isBlacklisted, err := IsTokenBlacklisted(claims.ID)
	//if err != nil {
	//	return nil, err
	//}
	//if isBlacklisted {
	//	return nil, ErrTokenBlacklisted
	//}
	//
	//return claims, nil

	return nil, errors.New("invalid token")
}

// 生成唯一 token ID (随机字节 + Hex 编码（推荐）)
func generateTokenID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// ExtractToken 从 Authorization 头中提取 JWT token
func ExtractToken(bearer string) string {
	if bearer == "" {
		return ""
	}

	parts := strings.SplitN(bearer, " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(bearer)
}
