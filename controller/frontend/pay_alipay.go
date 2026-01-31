package frontend

import (
	"bookstore/common"
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateAliPayOrder 创建支付宝扫码支付订单
func CreateAliPayOrder(w http.ResponseWriter, r *http.Request) {
	//1.获取订单 ID
	vars := mux.Vars(r)
	orderIDStr := vars["orderID"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "订单ID无效", nil)
		return
	}

	// 2.初始化支付宝支付服务
	paySrv := service.NewAliPayService()

	//3.创建支付订单(调用 service/payment_service.go)
	qrURL, err := paySrv.CreateAliPayOrder(orderID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	//4. 返回二维码 URL
	common.WriteJSON(w, http.StatusOK, "成功", map[string]any{
		"qr_url": qrURL,
	})

}
