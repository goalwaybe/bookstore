package frontend

import (
	"bookstore/common"
	"bookstore/dao"
	"bookstore/service"
	"bookstore/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
下单，查看订单详情
*/

func BuyNow(w http.ResponseWriter, r *http.Request) {
	//1.用户鉴权
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, http.StatusUnauthorized, "未登录或缺少用户ID", nil)
		return
	}

	//2.解析请求
	var req struct {
		BookID    int    `json:"book_id"`
		Count     int    `json:"count"`
		AddressID int    `json:"address_id"`
		Remark    string `json:"remark"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, "参数解析失败", nil)
		return
	}

	if req.BookID <= 0 || req.Count <= 0 || req.AddressID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "缺少必要参数", nil)
		return
	}

	//3.查商品
	book, err := dao.GetBookDetail(req.BookID)
	if err != nil || book == nil {
		common.WriteJSON(w, http.StatusBadRequest, "商品不存在", nil)
		return
	}

	//4.预扣库存 (Redis)
	err = dao.PreDeductStock(req.BookID, req.Count)
	if err != nil {
		common.WriteJSON(w, http.StatusBadRequest, "库存不足", nil)
		return
	}

	//5.计算金额
	totalAmount := float64(req.Count) * book.Price
	totalCount := req.Count

	//5.创建订单
	orderID, orderNo, err := dao.CreateOrder(
		userID,
		req.AddressID,
		totalAmount,
		totalCount,
		req.Remark,
	)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "创建订单失败", nil)
		return
	}

	//6.写入订单明细表
	err = dao.CreateOrderItem(orderID, req.BookID, req.Count, book.Price, orderNo)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "创建订单商品失败", nil)
		return
	}

	// 7.返回订单ID
	common.WriteJSON(w, http.StatusOK, "下单成功", map[string]interface{}{
		"order_id": orderID,
	})

}

// ConfirmOrder 提交订单(从购物车生成订单)
func ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	// -----------------------------------------------
	// 1.从中间件解析用户ID
	// -----------------------------------------------
	userID := utils.GetUserID(r)

	if userID == 0 {
		common.WriteJSON(w, http.StatusBadRequest, "未登录或缺少用户ID", nil)
		return
	}

	// -----------------------------------------------
	// 2.接受请求参数(如地址，备注)
	// -----------------------------------------------
	var req struct {
		AddressID int    `json:"address_id"`
		Remark    string `json:"remark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, "参数解析失败", nil)
		return
	}

	if req.AddressID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "请选择有效的收获地址", nil)
		return
	}

	// -----------------------------------------------
	// 3.获取购物车
	// -----------------------------------------------
	cart, err := dao.GetCartByUserID(userID)
	if err != nil || cart == nil {
		common.WriteJSON(w, http.StatusBadRequest, "购物车不存在", nil)
		return
	}

	//获取购物项
	items, err := dao.GetCartItems(cart.ID)
	if err != nil || len(items) == 0 {
		common.WriteJSON(w, http.StatusBadRequest, "购物车为空无法下单", nil)
		return
	}

	// -----------------------------------------------
	// 4.计算订单金额
	// -----------------------------------------------

	var totalAmount float64 = 0
	var totalCount int = 0
	for _, item := range items {
		book, err := dao.GetBookDetail(item.BookID)
		if err != nil || book == nil {
			common.WriteJSON(w, http.StatusBadRequest, "商品异常信息"+strconv.Itoa(item.BookID), nil)
			return
		}
		totalAmount += float64(item.Count) * book.Price
		totalCount += item.Count
	}

	// -----------------------------------------------
	// 5.创建订单主表
	// -----------------------------------------------
	orderID, orderNo, err := dao.CreateOrder(
		userID,
		req.AddressID,
		totalAmount,
		totalCount,
		req.Remark,
	)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "订单创建失败", nil)
		return
	}

	// -----------------------------------------------
	// 6. 订单商品表，逐个插入订单商品
	// -----------------------------------------------
	for _, item := range items {

		//预扣库存
		err = dao.PreDeductStock(item.BookID, item.Count)
		if err != nil {
			common.WriteJSON(w, http.StatusBadRequest, "商品库存不足", nil)
			return
		}

		err = dao.CreateOrderItem(
			orderID,
			item.BookID,
			item.Count,
			item.Price,
			orderNo,
		)
		if err != nil {
			common.WriteJSON(w, http.StatusInternalServerError, "创建订单商品失败", nil)
			return
		}
	}

	// -----------------------------------------------
	// 7.清空购物车
	// -----------------------------------------------
	if err := dao.ClearCart(cart.ID); err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "订单提交成功，但清空购物车失败", nil)
		return
	}

	// -----------------------------------------------
	// 8.返回订单 ID (前端跳转订单详情)
	// -----------------------------------------------
	common.WriteJSON(w, http.StatusOK, "订单提交成功",
		map[string]interface{}{
			"order_id": orderID,
		})
}

// GetOrderList 获取用户订单列表
func GetOrderList(w http.ResponseWriter, r *http.Request) {

	userID := utils.GetUserID(r)

	if userID == 0 {
		common.WriteJSON(w, http.StatusUnauthorized, "未登录或缺少用户ID", nil)
		return
	}

	//可选分页参数
	page := 1
	pageSize := 20
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	//查询订单列表
	orders, err := dao.GetOrdersByUserID(userID, page, pageSize)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "查询订单失败", nil)
		return
	}

	//返回结果
	common.WriteJSON(w, http.StatusOK, "查询成功", orders)
}

// GetOrderDetail 查看订单详情
func GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)

	if userID == 0 {
		common.WriteJSON(w, http.StatusBadRequest, "未登录或缺少用户ID", nil)
		return
	}

	vars := mux.Vars(r)
	OrderIdStr := vars["orderID"]
	orderID, err := strconv.Atoi(OrderIdStr)
	if err != nil || orderID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "订单ID无效", nil)
		return
	}

	order, err := dao.GetOrderByID(orderID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "查询订单失败", nil)
		return
	}

	if order == nil || order.UserID != userID {
		common.WriteJSON(w, http.StatusForbidden, "无权查看此订单", nil)
		return
	}

	common.WriteJSON(w, http.StatusOK, "查询成功", order)
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)

	if userID == 0 {
		common.WriteJSON(w, http.StatusBadRequest, "未登录或缺少用户ID", nil)
		return
	}

	//2.获取订单ID
	vars := mux.Vars(r)
	orderIDStr := vars["orderID"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "订单ID无效", nil)
		return
	}

	//3.查询订单是否存在且属于当前用户
	order, err := dao.GetOrderByID(orderID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "查询订单失败", nil)
		return
	}

	if order == nil || order.UserID != userID {
		common.WriteJSON(w, http.StatusForbidden, "无权取消此订单", nil)
		return
	}

	//4.检查订单状态，只有待支付或未发货订单可取消
	if order.State != 0 && order.State != 1 {
		common.WriteJSON(w, http.StatusBadRequest, "订单不可取消", nil)
		return
	}

	//5.执行取消操作
	err = dao.CancelOrderByID(orderID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "取消订单失败", nil)
		return
	}

	//取消订单返还库存
	items, _ := dao.GetOrderItemsByOrderID(orderID)
	for _, it := range items {
		dao.RollbackStock(it.BookID, it.Count)
	}

	common.WriteJSON(w, http.StatusOK, "订单已取消", nil)
}

// PayOrder 支付订单(模拟支付)
func PayOrder(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
	//1.获取用户ID
	userID := utils.GetUserID(r)

	if userID == 0 {
		common.WriteJSON(w, http.StatusBadRequest, "未登录或缺少用户ID", nil)
		return
	}

	//2.获取订单ID
	vars := mux.Vars(r)
	orderIDStr := vars["orderID"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "订单ID无效", nil)
		return
	}

	//3.查询订单信息
	order, err := dao.GetOrderByID(orderID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "查询订单失败", nil)
		return
	}

	if order == nil || order.UserID != userID {
		common.WriteJSON(w, http.StatusForbidden, "无权操作此订单", nil)
		return
	}

	// 4.订单状态必须是 “待支付”
	if order.State != 0 {
		common.WriteJSON(w, http.StatusBadRequest, "订单当前状态不可支付", nil)
		return
	}

	//5.更新订单状态为已支付
	err = dao.PayOrderByIDTx(tx, orderID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "订单支付失败", nil)
		return
	}

	//支付完成扣库存
	items, _ := dao.GetOrderItemsByOrderID(orderID)
	for _, it := range items {
		err = dao.ConfirmStockTx(tx, it.BookID, it.Count)
		if err != nil {
			common.WriteJSON(w, http.StatusInternalServerError, "库存写入失败", nil)
			return
		}
	}

	common.WriteJSON(w, http.StatusOK, "订单支付成功", nil)
}

// CheckPayStatus 前端轮询查询订单状态
func GetAliPayStatus(w http.ResponseWriter, r *http.Request) {
	orderID, _ := strconv.Atoi(r.URL.Query().Get("order_id"))
	// 2.初始化支付宝支付服务
	paySrv := service.NewAliPayService()
	status, err := paySrv.QueryAllPayStatus(orderID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(status))
}
