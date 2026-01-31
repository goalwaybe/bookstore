package service

import (
	"bookstore/config"
	"bookstore/dao"
	"context"
	"fmt"
	"net/http"

	"github.com/smartwalle/alipay/v3"
)

type PaymentService struct {
	client *alipay.Client
}

// NewAliPayService 初始化支付宝客户端
func NewAliPayService() *PaymentService {
	//创建支付宝客户端
	client, err := alipay.New(
		config.AliPayConfig.AppID,      // 应用 AppID
		config.AliPayConfig.PrivateKey, // 应用私钥
		false,                          // 沙箱环境 false
	)

	if err != nil {
		panic(fmt.Sprintf("初始化支付宝客户端失败:%v", err))
	}

	//加载支付宝公钥，用于回调验签
	if err := client.LoadAliPayPublicKey(config.AliPayConfig.AliPublicKey); err != nil {
		panic(fmt.Sprintf("加载支付宝公钥失败:%v", err))
	}

	return &PaymentService{client: client}
}

// CreateAliPayOrder 创建支付宝扫码支付订单，返回二维码 URL
func (s *PaymentService) CreateAliPayOrder(orderID int) (string, error) {
	//1.查询订单信息
	order, err := dao.GetOrderByID(orderID)
	if err != nil {
		return "", fmt.Errorf("查询订单失败:%v", err)
	}
	if order == nil {
		return "", fmt.Errorf("订单不存在")
	}

	//2.订单状态必须是待支付(0)
	if order.State != 0 {
		return "", fmt.Errorf("订单不是待支付状态")
	}

	// 2. 构造扫码支付预下单参数
	//p := alipay.TradePreCreate{
	//	Trade: alipay.Trade{
	//		OutTradeNo:  "ORD_TEST_123456",
	//		TotalAmount: "0.01",
	//		Subject:     "测试商品",
	//		ProductCode: "FACE_TO_FACE_PAYMENT",
	//		NotifyURL:   "http://你的回调地址/notify",
	//	},
	//	DiscountableAmount: "0.01",
	//	GoodsDetail: []*alipay.GoodsDetailItem{
	//		{
	//			GoodsId:   "1001",
	//			GoodsName: "测试商品A",
	//			Quantity:  1,
	//			Price:     "0.01",
	//		},
	//	},
	//}

	//3.构造支付宝预下单请求
	var p = alipay.TradePreCreate{
		Trade: alipay.Trade{
			OutTradeNo:  order.OrderNo,
			Subject:     "图书商城订单支付",
			ProductCode: "FACE_TO_FACE_PAYMENT", // 必填！
			TotalAmount: fmt.Sprintf("%.2f", order.TotalAmount),
			NotifyURL:   config.AliPayConfig.NotifyURL,
			ReturnURL:   config.AliPayConfig.ReturnURL,
		},
		//可选字段，根据需要设置
		DiscountableAmount: "0.00",
		//GoodsDetail: []*alipay.GoodsDetailItem{...},
	}

	//4.调用预下单接口
	ctx := context.Background()
	resp, err := s.client.TradePreCreate(ctx, p)
	if err != nil {
		fmt.Println("支付宝预订单失败:", err)
		return "", nil
	}

	if resp.QRCode == "" {
		return "", fmt.Errorf("支付宝返回二维码为空")
	}

	// 5.返回二维码 URL
	return resp.QRCode, nil
}

// QueryAllPayStatus 查询支付宝订单支付状态
func (s *PaymentService) QueryAllPayStatus(orderID int) (string, error) {
	order, err := dao.GetOrderByID(orderID)
	if err != nil {
	}
	if order == nil {
	}

	//构造查询请求
	outTradeNo := order.OrderNo
	tradeQuery := alipay.TradeQuery{}
	tradeQuery.OutTradeNo = outTradeNo
	ctx := context.Background()

	//构造查询请求
	resp, err := s.client.TradeQuery(ctx, tradeQuery)

	if err != nil {
		return "", fmt.Errorf("支付宝查询失败: %v", err)
	}

	//返回交易状态码（TRADE_SUCCESS、WAIT_BUYER_PAY、TRADE_CLOSED 等）
	return string(resp.TradeStatus), nil
}

// VerifyNotify 验证回调
func (s *PaymentService) VerifyNotify(r *http.Request) (bool, error) {
	//解析表单数据
	if err := r.ParseForm(); err != nil {
		return false, err
	}

	//调用 SDK 验证
	notify, err := s.client.DecodeNotification(r.Form)
	if err != nil {
		return false, err
	}
	//成功返回 true
	return notify != nil, nil
}
