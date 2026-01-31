package frontend

import (
	"bookstore/config"
	"bookstore/dao"
	"bookstore/model"
	"bookstore/service"
	"fmt"
	"net/http"
	"time"
)

func AliPayNotify(w http.ResponseWriter, r *http.Request) {
	//先解析 form
	if err := r.ParseForm(); err != nil {
		fmt.Printf("解析回调参数失败:", err)
		w.Write([]byte("fail"))
		return
	}

	// 1.验证签名
	paySrv := service.NewAliPayService()
	ok, err := paySrv.VerifyNotify(r)
	if err != nil || !ok {
		fmt.Println("支付宝回调验签失败:", err)
		w.Write([]byte("fail"))
		return
	}

	//2.解析订单号(支付宝回传)
	outTradeNo := r.FormValue("out_trade_no")
	if outTradeNo == "" {
		w.Write([]byte("fail"))
	}

	//3.查询订单 (直接使用 order_no 查询)
	order, err := dao.GetOrderByOrderNo(outTradeNo)
	if err != nil || order == nil {
		w.Write([]byte("fail"))
		return
	}

	//防止重复处理
	if order.State == 1 {
		_, werr := w.Write([]byte("success"))
		if werr != nil {
			return
		}
		return
	}

	// -----------------------------------------------
	// 启动事务
	// -----------------------------------------------
	tx, err := config.DB.Begin()
	if err != nil {
		_, err := w.Write([]byte("fail"))
		if err != nil {
			fmt.Printf("开启事务失败: %v\n", err)
			_, werr := w.Write([]byte("fail"))
			if werr != nil {
				return
			}
			return
		}
		return
	}

	//4.更新订单状态为已支付
	err = dao.PayOrderByIDTx(tx, order.ID)
	if err != nil {
		_, err := w.Write([]byte("fail"))
		if err != nil {
			fmt.Printf("更新订单失败:%v\n", err)
			_, werr := w.Write([]byte("fail"))
			if werr != nil {
				return
			}
			return
		}
		return
	}

	//5.扣库存
	items, _ := dao.GetOrderItemsByOrderID(order.ID)
	for _, it := range items {
		err = dao.ConfirmStockTx(tx, it.BookID, it.Count)
		if err != nil {
			tx.Rollback()
			fmt.Printf("扣减库存失败:%v, 商品ID:%d\n", err, it.BookID)
			w.Write([]byte("fail"))
			return
		}
	}

	//7.插入支付记录(必须在事务里)
	pay := &model.Payment{
		OrderNo:   order.OrderNo,
		OrderID:   order.ID,
		PayMethod: "alipay",
		PayTime:   time.Now(),
		Amount:    order.TotalAmount,
		Status:    1,
	}
	if err := dao.InsertPaymentTx(tx, pay); err != nil {
		tx.Rollback()
		fmt.Printf("插入支付记录失败:%v\n", err)
		w.Write([]byte("fail"))
		return
	}

	//8.提交事务
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Printf("提交事务失败:%v\n", err)
		w.Write([]byte("fail"))
		return
	}

	// 9. 记录成功日志
	fmt.Printf("支付宝支付成功: 订单号 %s, 金额 %.2f\n", order.OrderNo, order.TotalAmount)

	//6.回复支付宝 success
	w.Write([]byte("success"))
}
