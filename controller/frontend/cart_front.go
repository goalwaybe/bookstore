package frontend

import (
	"bookstore/common"
	"bookstore/dao"
	"bookstore/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
购物车(添加，删除，查看)
*/

// GetCart 获取购物车
func GetCart(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, http.StatusBadRequest, "未登录或token错误", nil)
		return
	}

	cart, err := dao.GetCartByUserID(userID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "获取购物车失败", nil)
		return
	}
	if cart == nil {
		common.WriteJSON(w, http.StatusOK, "购物车为空", nil)
		return
	}
	common.WriteJSON(w, http.StatusOK, "获取购物车成功", cart)
}

// AddToCart 添加商品
func AddToCart(w http.ResponseWriter, r *http.Request) {

	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, http.StatusUnauthorized, "用户未登录", nil)
		return
	}

	//定义请求的结构体
	type AddToCartRequest struct {
		BookID int `json:"book_id"`
	}

	var req AddToCartRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, "请求参数解析失败:"+err.Error(), nil)
		return
	}

	//参数验证
	if req.BookID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "商品ID无效", nil)
		return
	}

	//3.查询书籍真实价格
	book, err := dao.GetBookDetail(req.BookID)
	if err != nil || book == nil {
		common.WriteJSON(w, http.StatusBadRequest, "商品不存在", nil)
		return
	}
	price := book.Price

	cart, err := dao.GetCartByUserID(userID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "获取购物车失败", nil)
		return
	}
	var cartID int
	if cart == nil {
		newID, err := dao.CreateCart(userID)
		if err != nil {
			common.WriteJSON(w, http.StatusInternalServerError, "创建购物车失败", nil)
			return
		}
		cartID = int(newID)
	} else {
		cartID = cart.ID
	}

	item, err := dao.GetCartItem(cartID, req.BookID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "查询购物车商品失败", nil)
		return
	}

	if item == nil {
		err = dao.AddCartItem(cartID, req.BookID, price, 1)
	} else {
		err = dao.UpdateCartItem(item.ID, item.Count+1, price, cartID)
	}

	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "添加商品失败", nil)
		return
	}
	common.WriteJSON(w, http.StatusOK, "添加成功", nil)
}

// UpdateCartiItem 修改数量
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {

	//从 Header 获取用户ID
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, http.StatusUnauthorized, "未登录或缺少用户ID", nil)
		return
	}

	// 从URL 路径获取 bookID
	vars := mux.Vars(r)
	bookIDStr := vars["bookID"]
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "商品ID无效", nil)
		return
	}

	// 3.解析参数
	var req struct {
		Count int `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSON(w, http.StatusBadRequest, "参数解析失败", nil)
		return
	}

	if req.Count <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "数量必须大于0", nil)
		return
	}

	//查询购物车
	cart, err := dao.GetCartByUserID(userID)
	if err != nil || cart == nil {
		common.WriteJSON(w, http.StatusBadRequest, "购物车不存在", nil)
		return
	}

	//查询购物项
	item, err := dao.GetCartItem(cart.ID, bookID)
	if err != nil || item == nil {
		common.WriteJSON(w, http.StatusBadRequest, "商品不存在", nil)
		return
	}

	//获取真实价格
	book, err := dao.GetBookDetail(bookID)
	if err != nil || book == nil {
		common.WriteJSON(w, http.StatusBadRequest, "商品不存在", nil)
		return
	}

	//更新数量
	err = dao.UpdateCartItem(item.ID, req.Count, book.Price, cart.ID)

	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "修改商品数量失败", nil)
		return
	}
	//返回成功
	common.WriteJSON(w, http.StatusOK, "修改成功", nil)
}

// DeleteCartItem 删除商品
func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	// 从Header 获取用户ID
	userID := utils.GetUserID(r)
	if userID == 0 {
		common.WriteJSON(w, http.StatusUnauthorized, "未登录或缺少用户ID", nil)
		return
	}

	//从 URL 路径参数获取 bookID
	vars := mux.Vars(r)
	bookIDStr := vars["bookID"]
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		common.WriteJSON(w, http.StatusBadRequest, "商品ID无效", nil)
		return
	}

	cart, err := dao.GetCartByUserID(userID)
	if err != nil || cart == nil {
		common.WriteJSON(w, http.StatusBadRequest, "购物车不存在", nil)
		return
	}

	err = dao.DeleteCartItem(cart.ID, bookID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "删除商品失败", nil)
		return
	}
	common.WriteJSON(w, http.StatusOK, "删除成功", nil)
}

// ClearCart 清空购物车
func ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)

	if userID == 0 {
		common.WriteJSON(w, http.StatusUnauthorized, "未登录或缺少用户ID", nil)
		return
	}

	cart, err := dao.GetCartByUserID(userID)
	if err != nil || cart == nil {
		common.WriteJSON(w, http.StatusBadRequest, "购物车不存在", nil)
		return
	}

	err = dao.ClearCart(cart.ID)
	if err != nil {
		common.WriteJSON(w, http.StatusInternalServerError, "清空购物车", nil)
		return
	}
	common.WriteJSON(w, http.StatusOK, "购物车已清空", nil)
}
