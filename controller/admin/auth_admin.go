package admin

import (
	"bookstore/common"
	"bookstore/dao"
	"bookstore/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

/*
管理员登录 / 登出 /获取信息
*/

// AdminLoginHandler 管理员登录接口
func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		common.WriteJSON(w, 405, "Method Not Allowed", nil)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		common.WriteJSON(w, 400, "Invalid JSON", nil)
		return
	}

	// ---------- 分布式锁防并发登录 ----------
	lockKey := "lock:admin_login:" + req.Username
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

	admin, err := dao.GetAdminByUsername(req.Username)
	if err != nil {
		common.WriteJSON(w, 401, "用户不存在", nil)
		return
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		common.WriteJSON(w, 401, "密码错误", nil)
		return
	}

	// ---------- 单点登录处理 ----------
	oldTokenID, _ := dao.GetCache("admin_last_token:" + req.Username)
	if oldTokenID != "" {
		// 将旧 token 加入黑名单，有效期 2 小时（与 token 一致）
		_ = dao.BlacklistToken(oldTokenID, 2*time.Hour)
	}

	// 生成 Token (有效期 2 小时)
	token, claims, err := utils.GenerateToken(admin.ID, admin.Username, "admin", 2)
	if err != nil {
		common.WriteJSON(w, 500, "生成 Token 失败", nil)
		return
	}

	// 缓存当前 token ID，用于单点登录
	_ = dao.SetCache("admin_last_token:"+req.Username, claims.ID, 2*time.Hour)

	common.WriteJSON(w, 200, "登录成功", map[string]interface{}{
		"token":    token,
		"username": admin.Username,
		"role":     admin.Role,
	})
}

// AdminAddHandler 创建管理员
func AdminAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(common.JSONResponse{Code: 405, Msg: "Method Not Allowde"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.JSONResponse{Code: 400, Msg: "Invalid JSON"})
		return
	}

	if req.Username == "" || req.Password == "" {
		json.NewEncoder(w).Encode(common.JSONResponse{Code: 400, Msg: "用户名或密码不能为空"})
		return
	}

	// 调用 dao 层创建管理员
	err := dao.CreateAdmin(req.Username, req.Password)
	if err != nil {
		json.NewEncoder(w).Encode(common.JSONResponse{Code: 500, Msg: err.Error()})
		return
	}

	//返回成功
	json.NewEncoder(w).Encode(common.JSONResponse{
		Code: 200,
		Msg:  "管理员添加成功",
		Data: map[string]interface{}{
			"username": req.Username,
		},
	})

}

// AdminInfoHandler 获取当前登录的管理员信息
func AdminInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, 401, "未登录或 token 错误", nil)
		return
	}

	//查询数据库获取管理员信息
	admin, err := dao.GetAdminByID(userID)
	if err != nil {
		common.WriteJSON(w, 404, "管理员不存在", nil)
		return
	}

	// 返回管理员信息
	common.WriteJSON(w, 200, "获取成功", map[string]interface{}{
		"id":       admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
	})
}

// AdminLogoutHandler 登出接口
func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	//从请求头获取用户ID
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, 401, "未登录或token 错误", nil)
		return
	}

	// ---------- 从请求头中获取JWT并加入黑名单 ----------
	if err := dao.RevokeTokenFromRequest(r); err != nil {
		common.WriteJSON(w, 500, err.Error(), nil)
		return
	}

	common.WriteJSON(w, http.StatusOK, "注销成功", nil)
}

func AdminListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	keyword := r.URL.Query().Get("keyword")

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 {
		pageSize = 10
	}

	admins, total, err := dao.GetAdminList(keyword, page, pageSize)
	if err != nil {
		common.WriteJSON(w, 500, "系统错误:"+err.Error(), nil)
		return
	}
	common.WriteJSON(w, 200, "获取成功", map[string]interface{}{
		"list":      admins,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"totalPage": (total + pageSize - 1) / pageSize,
	})

}
