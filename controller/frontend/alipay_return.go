package frontend

import (
	"bookstore/dao"
	"fmt"
	"net/http"
)

// AliPayReturn 用户支付完成后的前端跳转
func AliPayReturn(w http.ResponseWriter, r *http.Request) {
	outTradeNo := r.URL.Query().Get("out_trade_no")

	order, err := dao.GetOrderByOrderNo(outTradeNo)
	if order == nil || err != nil {
		fmt.Fprintf(w, "订单不存在或查询失败!")
		return
	}

	// 根据订单状态显示支付结果
	if order.State == 1 {
		fmt.Fprintf(w, "支付处理中，请稍后刷新!订单号 : %s", outTradeNo)
	} else {
		fmt.Fprintf(w, "支付失败或订单异常! 订单号: %s", outTradeNo)
	}
}
