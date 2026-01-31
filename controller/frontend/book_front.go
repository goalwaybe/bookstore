package frontend

import (
	"bookstore/common"
	"bookstore/dao"
	"net/http"
	"strconv"
	"strings"
)

/*
图书浏览，详情，搜索
*/

// GetBookList 获取图书列表
func GetBookList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSON(w, 405, "只允许 GET 请求", nil)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	categoryStr := r.URL.Query().Get("category_id")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	categoryID, _ := strconv.Atoi(categoryStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	books, err := dao.GetBookList(categoryID, offset, limit)
	if err != nil {
		common.WriteJSON(w, 500, "数据库查询失败"+err.Error(), nil)
		return
	}

	common.WriteJSON(w, 200, "success", map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"category_id": categoryID,
		"list":        books,
	})

}

// GetBookDetai 获取图书详情
func GetBookDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSON(w, 405, "只允许 GET 请求", nil)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		common.WriteJSON(w, 400, "参数错误,缺少图书ID", nil)
		return
	}

	book, err := dao.GetBookDetail(id)
	if err != nil {
		common.WriteJSON(w, 500, "数据库查询失败："+err.Error(), nil)
		return
	}

	if book == nil || book.ID == 0 {
		common.WriteJSON(w, 404, "图书不存在", nil)
		return
	}
	common.WriteJSON(w, 200, "success", book)
}

// SearchBooks 搜索图书 (支持模糊搜索)
func SearchBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSON(w, 405, "只允许 GET 请求", nil)
		return
	}
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		common.WriteJSON(w, 400, "搜索关键词不能为空", nil)
		return
	}

	books, err := dao.SearchBooks(query)
	if err != nil {
		common.WriteJSON(w, 500, "数据库查询失败", nil)
		return
	}
	common.WriteJSON(w, 200, "success", map[string]interface{}{
		"keyword": query,
		"count":   len(books),
		"list":    books,
	})
}

// GetCategoryList 获取分类
func GetCategoryList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteJSON(w, 405, "只允许 GET 请求", nil)
		return
	}

	categories, err := dao.GetCategoryList()
	if err != nil {
		common.WriteJSON(w, 500, "数据库查询失败"+err.Error(), nil)
		return
	}
	if len(categories) == 0 {
		common.WriteJSON(w, 404, "暂无分类数据", nil)
		return
	}
	common.WriteJSON(w, 200, "success", categories)
}
