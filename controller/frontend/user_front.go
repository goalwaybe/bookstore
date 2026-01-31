package frontend

import (
	"bookstore/common"
	"bookstore/dao"
	"bookstore/service"
	"bookstore/utils"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

/*
用户注册，登录，信息修改
*/

// RegisterRequest 定义结构体接收 JSON 数据
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Register 注册用户
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, 400, "请求数据格式错误", nil)
		return
	}

	if req.Username == "" || req.Password == "" || req.Email == "" {
		common.WriteJSON(w, 400, "用户名、密码或邮箱不能为空", nil)
		return
	}

	// 检查用户名是否存在
	user, err := dao.CheckUserName(req.Username)
	if err != nil {
		common.WriteJSON(w, 500, "数据库查询失败", nil)
		return
	}

	if user != nil && user.ID > 0 {
		common.WriteJSON(w, 400, "用户名已存在", nil)
		return
	}

	//检查邮箱是否存在
	userByEmail, err := dao.CheckEmail(req.Email)
	if err != nil {
		common.WriteJSON(w, 500, "数据库查询失败", nil)
		return
	}

	if userByEmail != nil && userByEmail.ID > 0 {
		common.WriteJSON(w, 400, "邮箱已被注册", nil)
		return
	}

	//保存用户
	err = dao.SaveUser(req.Username, req.Password, req.Email)
	if err != nil {
		common.WriteJSON(w, 500, "注册失败", nil)
		return
	}

	common.WriteJSON(w, 200, "注册成功", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// 解析 JSON 请求
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, 400, "请求数据格式错误", nil)
		return
	}

	if req.Username == "" || req.Password == "" {
		common.WriteJSON(w, 400, "用户名或密码不为为空", nil)
		return
	}

	// ---------- 分布式锁防止并发登录 ----------
	lockKey := "lock:user_login:" + req.Username
	gotLock, err := dao.TryLock(lockKey, 2*time.Second)
	if err != nil {
		common.WriteJSON(w, 500, "系统错误，请稍后重试", nil)
		return
	}
	if !gotLock {
		common.WriteJSON(w, 429, "登录请求过于频繁，请稍后重试", nil)
		return
	}
	defer dao.Unlock(lockKey)

	// 检查用户名密码
	user, err := dao.CheckUserNameAndPassword(req.Username, req.Password)
	if err != nil {
		common.WriteJSON(w, 500, "数据库查询错误", nil)
		return
	}

	if user.ID == 0 {
		common.WriteJSON(w, 401, "用户名或密码不正确", nil)
		return
	}

	// ---------- 单点登录处理 ----------
	oldTokenID, _ := dao.GetCache("user_last_token:" + req.Username)
	if oldTokenID != "" {
		//将旧token加入黑名单，有效期2小时（与token一致）
		_ = dao.BlacklistToken(oldTokenID, 2*time.Hour)
	}

	//生成 token
	token, claims, err := utils.GenerateToken(user.ID, user.Username, "user", 2)
	if err != nil {
		common.WriteJSON(w, 500, "生成 token 失败", nil)
		return
	}

	// ---------- 缓存当前 token ID ----------
	_ = dao.SetCache("user_last_token:"+req.Username, claims.ID, 2*time.Hour)

	//返回登录成功 JSON
	common.WriteJSON(w, 200, "登录成功", map[string]interface{}{
		"userID":   user.ID,
		"username": user.Username,
		"userType": "user",
		"token":    token,
	})
}

// Logout 用户登出
func Logout(w http.ResponseWriter, r *http.Request) {
	// 从请求头获取 Token
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		common.WriteJSON(w, http.StatusBadRequest, "未提供 Token", nil)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		common.WriteJSON(w, http.StatusBadRequest, "Authorization 格式错误", nil)
		return
	}

	tokenString := parts[1]

	// 解析 Token
	claims, err := utils.ParseToken(utils.ExtractToken(tokenString))
	if err != nil {
		common.WriteJSON(w, http.StatusUnauthorized, "Token 无效或已过期", nil)
		return
	}

	// 计算 Token 剩余有效期
	expireDuration := time.Until(claims.ExpiresAt.Time)
	if expireDuration <= 0 {
		common.WriteJSON(w, http.StatusUnauthorized, "Token 已过期", nil)
		return
	}

	// 将 Token 加入黑名单，防止继续使用
	if err := dao.BlacklistToken(claims.ID, expireDuration); err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "登出失败", nil)
		return
	}

	common.WriteJSON(w, http.StatusOK, "登出成功", nil)
}

// GetUserInfo 获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, 401, "未登录或 Token 错误", nil)
		return
	}

	//查询数据库
	user, err := dao.GetUserByID(userID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "数据库查询失败", nil)
		return
	}

	if user == nil {
		common.WriteJSON(w, http.StatusNotFound, "用户不存在", nil)
		return
	}
	common.WriteJSON(w, http.StatusOK, "success", user)
}

// UpdateProfile 修改个人资料
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, 401, "未登录或 Token 错误", nil)
		return
	}

	// 解析请求 JSON
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		common.WriteJSON(w, 400, "请求参数格式错误", nil)
		return
	}

	//调用 dao 层更新
	updatedUser, err := dao.UpdateProfileByID(userID, req.Username, req.Email)
	if err != nil {
		common.WriteJSON(w, 400, err.Error(), nil)
		return
	}

	if updatedUser == nil {
		common.WriteJSON(w, 404, "用户不存在", nil)
		return
	}

	common.WriteJSON(w, 200, "资料更新成功", updatedUser)
}

// UpdatePassword 修改密码
func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	// 从 Header 获取用户ID
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, 401, "未登录或 Token 错误", nil)
		return
	}

	// 解析前端 JSON 请求体
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, 400, "请求参数格式错误", nil)
		return
	}

	// 参数校验
	if req.OldPassword == "" || req.NewPassword == "" {
		common.WriteJSON(w, 400, "旧密码和新密码都不能为空", nil)
		return
	}

	if len(req.NewPassword) < 6 {
		common.WriteJSON(w, 400, "新密码长度不能少于6位", nil)
		return
	}

	// 调用 DAO 层逻辑更新密码
	err := dao.UpdatePasswordByID(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		common.WriteJSON(w, 400, err.Error(), nil)
		return
	}

	// ---------- 从请求头中获取 JWT 并加入黑名单  ----------
	//将当前Token 拉入黑名单，过期时间建议和Token一致
	if err := dao.RevokeTokenFromRequest(r); err != nil {
		common.WriteJSON(w, 500, err.Error(), nil)
		return
	}

	// 返回成功响应
	common.WriteJSON(w, 200, "密码修改成功，请重新登录", nil)
}

// SendRecoverCode 发送邮箱验证码
func SendRecoverCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		UserType string `json:"user_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, 400, "参数解析错误", nil)
		return
	}

	err := service.SendResetCodeService(req.Email)
	if err != nil {
		common.WriteJSON(w, 400, err.Error(), nil)
		return
	}
	common.WriteJSON(w, 200, "验证码已发送，请检查邮箱", nil)
}

// RecoverPassword 找回密码
func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"new_password"`
		UserType    string `json:"user_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, 400, "参数解析错误", nil)
		return
	}

	if req.Email == "" || req.Code == "" || req.NewPassword == "" {
		common.WriteJSON(w, 400, "邮箱、验证码、新密码都不能为空", nil)
		return
	}

	err := service.RecoverPasswordService("user", req.Email, req.Code, req.NewPassword)
	if err != nil {
		common.WriteJSON(w, 400, err.Error(), nil)
		return
	}
	common.WriteJSON(w, 200, "密码重置成功", nil)
}
