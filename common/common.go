package common

import (
	"encoding/json"
	"net/http"
)

/*
JSONResponse 前后台共享的通用控制器（如验证码、上传文件等）
*/
type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// WriteJSON 工具函数：标准响应
func WriteJSON(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSONResponse{Code: code, Msg: msg, Data: data})
}
