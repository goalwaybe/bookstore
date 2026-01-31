package router

import (
	"bookstore/controller/admin"
	"bookstore/controller/frontend"
	"net/http"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	// --------------------------------------------
	// 后台模块
	// --------------------------------------------

	// ---------- 管理员操作管理员自身管理模块 ----------
	//管理员登录
	r.HandleFunc("/admin/login", admin.AdminLoginHandler).Methods("POST")

	//管理员信息 应用中间件
	r.Handle("/admin/info", JWTAuthMiddleware(http.HandlerFunc(admin.AdminInfoHandler))).Methods("GET")

	//管理员退出
	r.Handle("/admin/logout", JWTAuthMiddleware(http.HandlerFunc(admin.AdminLogoutHandler))).Methods("POST")

	//新增管理员
	r.Handle("/admin/add", JWTAuthMiddleware(http.HandlerFunc(admin.AdminAddHandler))).Methods("POST")

	//管理员列表
	r.Handle("/admin/list", JWTAuthMiddleware(http.HandlerFunc(admin.AdminListHandler))).Methods("GET")

	//管理员更新
	//r.Handle("/admin/update/{adminID}",JWTAuthMiddleware(http.HandlerFunc(admin.AdminUpdateHandler))).Methods("PUT")

	//管理员删除
	//r.Handle("/admin/delete/{adminID}",JWTAuthMiddleware(http.HandlerFunc(admin.AdminDeleteHandler))).Methods("DELETE")

	// ---------- 管理员操作用户管理模块 ----------
	// /admin/user/list 				GET		获取用户列表(分页，搜索)
	// /admin/user/detail/{userID} 		GET		查看指定用户信息
	//  /admin/user/add 				POST	新增用户
	//  /admin/user/update/{userID}		PUT		更新用户信息（昵称、状态、权限等)
	//  /admin/user/delete/{userID}		DELETE	删除指定用户
	//	/admin/user/ban/{userID}		POST	封禁用户
	//	/admin/user/unban/{userID}		POST	解封用户

	// ---------- 管理员操作图书管理模块 ----------
	// 	/admin/book/list				GET		获取图书列表
	//	/admin/book/detail/{bookID}		GET		查看图书详情
	//	/admin/book/add					POST	添加图书
	//	/admin/book/update/{bookID}		PUT		更新图书信息（价格、库存、封面、描述等)
	//	/admin/book/delete/{bookID}		DELETE	删除图书
	//	/admin/book/batchUpdateStock	PUT		批量修改库存
	//	/admin/book/batchDelete			DELETE	批量删除图书

	// ---------- 管理员操作分类管理模块 ----------
	// /admin/category/list						GET			获取分类列表
	//	/admin/category/detail/{categoryID}		GET			查看分类详情
	//	/admin/category/add						POST		新增分类
	//	/admin/category/update/{categoryID}		PUT			更新分类
	//	/admin/category/delete/{categoryID} 	DELETE		删除分类

	// ---------- 管理员操作订单管理模块 ----------
	// /admin/order/list 	GET		获取订单列表 (分页、状态筛选)
	// /admin/order/detail/{orderID}   	GET  查看订单详情
	// /admin/order/cancel/{orderID}	POST  取消订单
	// /admin/order/refund/{orderID}	POST 发起退款
	//

	// ---------- 管理员操作购物车管理模块 ----------
	// /admin/cart/list/{userID}	GET		查看用户指定购物车

	// ---------- 管理员操作收货地址管理模块 ----------
	//  /admin/address/list/{userID}  GET		获取用户收货地址列表
	//	/admin/address/delete/{addressID}		DELETE		删除异常地址

	// ---------- 管理员操作支付与财务管理模块 ----------
	//  /admin/payment/list		GET		查看支付记录
	//  /admin/payment/detail/{orderID}		GET		查看订单支付详情
	//	/admin/payment/refund/{orderID}		POST		手动退款
	//	/admin/finance/report				GET			获取财务报表 (销售额、利润等）

	// -----------------------------------------------
	// 前台模块
	// -----------------------------------------------

	// ---------- 用户相关 ----------
	//用户注册
	r.HandleFunc("/user/register", frontend.Register).Methods("POST")

	//用户登录
	r.HandleFunc("/user/login", frontend.Login).Methods("POST")

	// 用户登出
	r.Handle("/user/logout", JWTAuthMiddleware(http.HandlerFunc(frontend.Logout))).Methods("POST")

	// 获取用户信息
	r.Handle("/user/info/", JWTAuthMiddleware(http.HandlerFunc(frontend.GetUserInfo))).Methods("GET")

	//用户修改密码
	r.Handle("/user/updatePassword", JWTAuthMiddleware(http.HandlerFunc(frontend.UpdatePassword))).Methods("PUT")

	//用户发送验证码
	r.HandleFunc("/user/sendrecovercode", frontend.SendRecoverCode).Methods("POST")

	//用户重置密码
	r.HandleFunc("/user/recoverpassword", frontend.RecoverPassword).Methods("POST")

	//用户注销账户，通过 JWT 获取 userID
	//r.Handle("/user/delete", JWTAuthMiddleware(http.HandlerFunc(frontend.DeleteUser))).Methods("DELETE")

	// ---------- 图书展示模块 ----------
	// 获取图书列表（无需登录，公开接口）
	r.HandleFunc("/book/list", frontend.GetBookList).Methods("GET") //获取图书列表

	//获取指定 bookID 的图书详情（无需登录，公开接口）
	r.HandleFunc("/book/detail/{bookID}", frontend.GetBookDetail).Methods("GET") //获取图书详情

	// 根据关键词搜索图书（无需登录，公开接口）
	r.HandleFunc("/book/search", frontend.SearchBooks).Methods("GET") //搜索图书

	// 获取图书分类列表（无需登录，公开接口）
	r.HandleFunc("/category/list", frontend.GetCategoryList).Methods("POST") //获取分类

	// ---------- 购物车模块 ----------
	//获取购物车详情
	r.Handle("/cart/get", JWTAuthMiddleware(http.HandlerFunc(frontend.GetCart))).Methods("GET")

	//添加商品
	r.Handle("/cart/add", JWTAuthMiddleware(http.HandlerFunc(frontend.AddToCart))).Methods("POST")

	//修改商品
	r.Handle("/cart/update/{bookID}", JWTAuthMiddleware(http.HandlerFunc(frontend.UpdateCartItem))).Methods("POST")

	//删除商品
	r.Handle("/cart/delete/{bookID}", JWTAuthMiddleware(http.HandlerFunc(frontend.DeleteCartItem))).Methods("DELETE")

	// 清空购物车
	r.Handle("/cart/clear", JWTAuthMiddleware(http.HandlerFunc(frontend.ClearCart))).Methods("POST")

	// ---------- 订单模块 ----------

	//直接购买生成订单
	r.Handle("/order/buy", JWTAuthMiddleware(http.HandlerFunc(frontend.BuyNow))).Methods("POST")

	//提交订单（从购物车生成订单）
	r.Handle("/order/confirm", JWTAuthMiddleware(http.HandlerFunc(frontend.ConfirmOrder))).Methods("POST")

	//订单列表
	r.Handle("/order/list", JWTAuthMiddleware(http.HandlerFunc(frontend.GetOrderList))).Methods("GET")

	//查看订单详情
	r.Handle("/order/detail/{orderID}", JWTAuthMiddleware(http.HandlerFunc(frontend.GetOrderDetail))).Methods("GET")

	//取消订单
	r.Handle("/order/cancel/{orderID}", JWTAuthMiddleware(http.HandlerFunc(frontend.CancelOrder))).Methods("POST")

	//支付订单（模拟支付）接入真实支付后废弃
	//r.Handle("/order/pay/{orderID}", JWTAuthMiddleware(http.HandlerFunc(frontend.PayOrder))).Methods("POST")

	// -----------------------------------------------
	// 支付模块(前台)
	// -----------------------------------------------
	//1.用户生成支付宝支付订单(返回二维码或支付URL)
	r.Handle("/pay/alipay/create/{orderID}", JWTAuthMiddleware(http.HandlerFunc(frontend.CreateAliPayOrder))).Methods("POST")

	//2.用户查询支付宝支付状态 ，查询订单查询失败还是成功
	r.Handle("/pay/alipy/status/{orderID}", JWTAuthMiddleware(http.HandlerFunc(frontend.GetAliPayStatus))).Methods("GET")

	//3.用户生成微信支付订单(JSAPI)
	//r.Handle("/pay/wechat/create/{orderID}", JWTAuthMiddleware(http.HandlerFunc(frontend.CreateWxPayOrder))).Methods("POST")

	//4.用户查询微信支付状态
	//r.Handle("/pay/wechat/status/{orderID}", JWTAuthMiddleware(http.HandlerFunc(frontend.GetWxPayStatus))).Methods("GET")

	// -----------------------------------------------
	// 支付宝回调
	// -----------------------------------------------

	//支付宝异步回调
	r.HandleFunc("/pay/alipay/notify", frontend.AliPayNotify).Methods("POST")

	//支付宝付款成功后返回订单号和跳转
	//r.HandleFunc("/pay/alipay/return", frontend.AliPayReturn).Methods("GET")

	//微信支付异步回调
	//r.HandleFunc("/pay/wechat/notify", frontend.WxPayNotify).Methods("POST")

	// -----------------------------------------------
	// 后台管理支付逻辑，查询订单，退款等
	// -----------------------------------------------

	//后台管理查询所有支付订单
	//r.Handle("/admin/pay/list", JWTAuthMiddleware(http.HandlerFunc(admin.ListPayOrder))).Methods("POST")
	//
	////后台人工退款(支付宝退款)
	//r.Handle("/admin/pay/refund/alipay/{orderID}", JWTAuthMiddleware(http.HandlerFunc(admin.RefundAlipy))).Methods("POST")
	//
	////后台人工退款(微信退款)
	//r.Handle("/admin/pay/refund/wechat/{orderID}", JWTAuthMiddleware(http.HandlerFunc(admin.RefundWxPay))).Methods("POST")

	// ---------- 收货地址模块 ----------

	// ---------- 支付接口(模拟) 考虑使用 支付宝沙箱模式----------

	return r
}
