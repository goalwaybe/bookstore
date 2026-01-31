package main

import (
	"bookstore/config"
	"bookstore/cron"
	"bookstore/router"
	"bookstore/service"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/smartwalle/alipay/v3"
)

func main() {
	//
	//config.InitRedis()
	//
	//err := dao.SetCache("book_123", "rrrr", 10*time.Minute)
	//if err != nil {
	//	log.Println("设置缓存失败:", err)
	//}
	//
	//val, err := dao.GetCache("book_123")
	//if err != nil {
	//	log.Println("获取缓存失败：", err)
	//} else {
	//	log.Println("缓存值：", val)
	//}

	////管理员登录
	//http.HandleFunc("/adminlogin", admin.AdminLoginHandler)
	//
	////管理员信息
	//http.HandleFunc("/admininfo", admin.AdminInfoHandler)
	//
	////管理员登出
	//http.HandleFunc("/adminlogout", admin.AdminLogoutHandler)
	//
	////创建管理员
	//http.HandleFunc("/adminAdd", admin.AdminAddHandler)
	//
	////添加图书到购物车
	//http.HandleFunc("/addBook2Cart", controller.AddBook2Cart)
	//
	////获取带分页和价格范围的图书
	//http.HandleFunc("/getPageBooksByPrice", controller.GetPageBooksByPrice)
	//
	////获取带分页的图书信息
	//http.HandleFunc("/getPageBooks", controller.GetPageBooks)
	//
	////注册删除图书接口
	//http.HandleFunc("/deleteBook", controller.DeleteBook)
	//
	////根据ID查询图书
	//http.HandleFunc("/getBookByID", controller.GetBookByIDHandler)
	//
	//// 新增或者更新书籍信息
	//http.HandleFunc("/updateOrAddBook", controller.UpdateOrAddBookHandler)
	//
	////获取所有分类
	//http.HandleFunc("/getAllCategoriesHandler", controller.GetAllCategoriesHandler)
	//
	////获取分类树形数据
	//http.HandleFunc("/getCategoryTreeHandler", controller.GetCategoryTreeHandler)
	//
	////查看用户是否存在, 校验用户
	//http.HandleFunc("/checkUserNameHandler", controller.CheckUserNameHandler)
	//
	////注册新用户
	//http.HandleFunc("/regist", controller.Regist)
	//
	////登录
	//http.HandleFunc("/login", controller.Login)
	//
	////注销
	//http.HandleFunc("/logout", controller.Logout)

	// ---------- 发送邮件测试 ----------
	// 邮件服务器配置
	//from := "wangkxin@foxmail.com" // 发件人邮箱
	//password := "mkjjkjkjkfdfdfd"  // ⚠️ 注意：这里不是邮箱密码，而是“授权码”
	//to := "846166458@qq.com"       // 收件人邮箱（可换成自己的另一个邮箱）
	//host := "smtp.qq.com"
	//port := 465
	//
	//// 邮件内容 (HTML)
	//subject := "测试邮件"
	//body := `
	//	<h2>Go 邮件发送测试</h2>
	//	<p>这是一封通过 <b>Go 语言</b> 发送的 HTML 邮件。</p>
	//`
	//
	//msg := []byte(fmt.Sprintf(
	//	"From: %s\r\n"+
	//		"To: %s\r\n"+
	//		"Subject: %s\r\n"+
	//		"MIME-Version: 1.0\r\n"+
	//		"Content-Type: text/html; charset=UTF-8\r\n\r\n%s",
	//	from, to, subject, body,
	//))
	//
	//// 连接 SMTP 服务器（TLS）
	//addr := fmt.Sprintf("%s:%d", host, port)
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: true,
	//	ServerName:         host,
	//}
	//
	//conn, err := tls.Dial("tcp", addr, tlsConfig)
	//if err != nil {
	//	fmt.Println("❌ TLS 连接失败:", err)
	//	return
	//}
	//
	//client, err := smtp.NewClient(conn, host)
	//if err != nil {
	//	fmt.Println("❌ 创建客户端失败:", err)
	//	return
	//}
	//defer client.Close()
	//
	//// 登录认证
	//auth := smtp.PlainAuth("", from, password, host)
	//if err = client.Auth(auth); err != nil {
	//	fmt.Println("❌ 登录失败:", err)
	//	return
	//}
	//
	//// 发件人与收件人
	//if err = client.Mail(from); err != nil {
	//	fmt.Println("❌ 设置发件人失败:", err)
	//	return
	//}
	//if err = client.Rcpt(to); err != nil {
	//	fmt.Println("❌ 设置收件人失败:", err)
	//	return
	//}
	//
	//// 写入邮件内容
	//writer, err := client.Data()
	//if err != nil {
	//	fmt.Println("❌ 发送数据失败:", err)
	//	return
	//}
	//_, err = writer.Write(msg)
	//if err != nil {
	//	fmt.Println("❌ 写入邮件内容失败:", err)
	//	return
	//}
	//err = writer.Close()
	//if err != nil {
	//	fmt.Println("❌ 关闭写入通道失败:", err)
	//	return
	//}
	//
	//client.Quit()
	//fmt.Println("✅ 邮件发送成功！")

	//TestAliPayOrder()

	// ---------- 新的配置 ----------

	// 1️⃣ 先加载 YAML 配置
	config.InitConfig()

	// 2️⃣ 初始化数据库
	config.InitDB()
	defer config.DB.Close()

	//初始化 Redis
	config.InitRedis()
	defer config.RedisClient.Close()

	// 初始化雪花节点
	config.InitSnowflake(1)

	//启动 Cron 任务 同步库存到 redis
	cron.InitCronJobs()

	//在这里执行一次库存批量同步
	stockService := &service.StockSyncService{}
	if err := stockService.SyncAllStock(); err != nil {
		log.Println("启动同步失败:", err)
	}

	// 初始化路由
	r := router.InitRoutes()

	port := ":8030"
	fmt.Println("Starting server at", port)

	err := http.ListenAndServe(port, r)

	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

// TestAliPayOrder 测试支付宝沙箱扫码支付
func TestAliPayOrder() {
	// 1. 初始化支付宝客户端（公钥模式）
	client, err := alipay.New(
		"9021000157684038", // AppID，沙箱账号生成的
		"MIIEowIBAAKCAQEAjyJZ+8t+i2xHQV5Kg4IP+bBnJ5wx1Vifv7r6l86qhsA1LqrVpLajTVmkgKQsR0a9nKSoo4HLdCl8B0K4PSVPJNUudwMFxFzVPhY7sTvEaoaGjWG/qbLM5IO4a5cAr4fjm6wFI18XwOcrlMAzBIZlLIN2p4iU2eIcePKkhsxypLSciVt0QM7fQoW+Tly7QPMbAtp3ev5PvCj/PLkRHHJsS8vFAw1QSTJJq6RNluZnZuhVMwXBbuPz6OhcMI/FEKJT3gm2pyC1P/iebfd6IMzN6SC27BI5w1MMC3/jBixZduOqgFIrmI9SZb6lifXSXHBspFbO3llwf5lMEqLH/b5SswIDAQABAoIBAGbh9zIRBEWtL86huwep/leoX41OLVe/2xvSh3uE1ngDQnbn3qcFjH312FOLLjSReYtVo83bZuV5SChj4dA5uBiujzaGWqfZXoHBM2jin/098ws+3qhY2APNVksngBurnoPv8sWI/abvwxipykIPZDvJxwkHGI9AgtzA9FfcLu+JuXL/tggUEF86tT8AgSuQ4uNKf4/6IgSWvw5wnYeIB6mbRY7beOoBQWCQXBOzAo/FQzVAE7GiY/KTLLVtQeYQz3P/4HM+J4rFK0nFGoCaFou7lk71y8mKMFGIK/tW9hVIdKA5UlY2h4xY94wT2hNrr8eYy3e/jy1Q/S6w2RjRAzkCgYEA2r8oB0zl9wy8zwGU7uRUg1eaw8rGPN5TAEXw6VMOl0EO7J6nXhWj7sunI+2Ce26rRBJ9JTcM8kzSHmOOpXb1GPj25HB1xeWtBhVng2EdR7KgphH9wL2JDqPkwmU/oGf8qHMzUjQLQbvoOmlyXZfU7X/FWYBGG/F+88f3oCYCL80CgYEAp4KqoTpOlTrnWUXkFWZ2AXA78/QTN6X4oEjDGf+r6Hdy0QW67PXCmuH6trkYag4j6TgfezoruHZfCL26do92JjOWpbMNoYVogux5h25uouQx6xy3ajIV0qBYIQGZI5JkHVsIyEbVbkcXT8fam3s/zkWqjelxBTfq0s6uhrEXDH8CgYA/A3gFelp4q3usajk8BBxwZYkYw84NdMIWPw+iNyHp4nzpHt751GQQAyyXxfCjnALMLkNMuCoLlqOky2spMTZzxixhLCBRLNOGAB58kzo7EDMFBAPERRU5WW8prL8Cd1IqOF1dXc6sVfQU95FRcs6MyBTSrogDvrKiiUVFJg3VNQKBgDfiN9qiASV4CUaXMoiRLj08nEO+4cpm7JNMswoxacCcWQmyx7HWK55QVbwv9B0tqn1b4+TQe1WG6B2abmKvadNE0rVlDt8cCsni6ohsJyzgxGaTpf0fyHgEVmJsjhJ3/D5u6Hcoanvn67cKDbxhWeAPDd9aSgJbrVO28DKvUekxAoGBAMeNrxyLO0sRPDdcWr/5/Mv5tnM/vS04aHDAGPsV3pl3o2xqN1EblLgQl37AY1dqk69h11VrBSwqSmIiYlXRm1bkBJYY9hZ6FWxSdJ1TzSy28K+FXuI3vnJ813YkFk5Uy8PYdH0mBXQKR070AMpmiwBc6L68eOdW2BZnKSCmrUp5", // 私钥（必须填）
		false,
	)
	if err != nil {
		fmt.Println("初始化支付宝客户端失败:", err)
		return
	}

	// 2. 加载支付宝公钥（回调验签必须）
	err = client.LoadAliPayPublicKey("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmv83o+GpX6rNQp55OQvOmCWOSXzN8w9QoYBXFPbMCKM68yw1j+k5/6xpZ+b/GB5b8vjy4Fb6hVuKJlXxt9CHy2RzUjUw1wr2CvrOjOInzPuZieaCLvYbGw/2hMCF7wOOFVBSJgLM6K5pPqIHsDddjxcL4LmFcBcRZmma5na158J1Tm65j9XA9n1vODmOxx5eIfuMCZdb2Wyiq2tKPRuCrUb2uSR3zqrtWb2klxmxpsX5dp+4EjuGBBeN8vf1V/PyDJxuu6AODnD+ZY/LvSszOQt2Pjc2i83bPRBAoSSjIJwFuygZo+3SbHckstkl2hP6oGJhmyOvibwrfCeS9n4udwIDAQAB")
	if err != nil {
		fmt.Println("加载支付宝公钥失败:", err)
		return
	}

	// 2. 构造扫码支付预下单参数
	p := alipay.TradePreCreate{
		Trade: alipay.Trade{
			OutTradeNo:  "ORD_TEST_123456",
			TotalAmount: "0.01",
			Subject:     "测试商品",
			ProductCode: "FACE_TO_FACE_PAYMENT",
			NotifyURL:   "http://你的回调地址/notify",
		},
		DiscountableAmount: "0.01",
		GoodsDetail: []*alipay.GoodsDetailItem{
			{
				GoodsId:   "1001",
				GoodsName: "测试商品A",
				Quantity:  1,
				Price:     "0.01",
			},
		},
	}

	// 3. 调用预下单接口
	ctx := context.Background() // 创建上下文
	resp, err := client.TradePreCreate(ctx, p)
	if err != nil {
		fmt.Println("支付宝预订单失败:", err)
		return
	}

	// 4. 输出二维码 URL
	if resp.QRCode == "" {
		fmt.Println("支付宝返回二维码为空")
		return
	}

	fmt.Println("预订单二维码 URL:", resp.QRCode)
}
