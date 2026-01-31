package dao

import (
	"bookstore/config"
	"bookstore/model"
	"database/sql"
)

// GetBookList 获取图书列表(带分类名)
func GetBookList(categoryID, offset, limit int) ([]model.Book, error) {
	queryBase := `
SELECT b.id,b.title,b.author,b.price,b.sales,b.stock,b.img_path,
b.category_id,IFNULL(c.name,'') AS category_name
FROM books AS b
LEFT JOIN categories AS c ON b.category_id = c.id
`
	var rows *sql.Rows
	var err error

	// 如果传入分类ID，则筛选分类；否则默认展示畅销/最新
	if categoryID > 0 {
		queryBase += "WHERE b.category_id = ? ORDER BY b.sales DESC LIMIT ?,?"
		rows, err = config.DB.Query(queryBase, categoryID, offset, limit)
	} else {
		queryBase += "ORDER BY b.sales DESC, b.id DESC LIMIT ?,?"
		rows, err = config.DB.Query(queryBase, offset, limit)
	}

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var books []model.Book
	for rows.Next() {
		var b model.Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.Price, &b.Sales, &b.Stock, &b.ImgPath, &b.CategoryID, &b.CategoryName)
		books = append(books, b)
	}
	return books, nil
}

// GetBookDetail 获取本图书详情
func GetBookDetail(id int) (*model.Book, error) {
	query := `
SELECT b.id,b.title,b.author,b.price,b.sales,b.stock,b.img_path,
b.category_id, IFNULL(c.name,'') AS category_name
FROM books AS b
LEFT JOIN categories AS c ON b.category_id = c.id
WHERE b.id = ?
`
	var b model.Book
	err := config.DB.QueryRow(query, id).Scan(&b.ID, &b.Title, &b.Author, &b.Price, &b.Sales, &b.Stock, &b.ImgPath, &b.CategoryID, &b.CategoryName)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// SearchBooks 模糊搜索图书
func SearchBooks(keyword string) ([]model.Book, error) {
	query := `
SELECT b.id,b.title,b.author,b.price,b.sales,b.stock,b.img_path,
b.category_id,IFNULL(c.name,'') AS category_name
FROM books AS b
LEFT JOIN categories AS c ON b.category_id = c.id
WHERE b.title LIKE ?
`
	rows, err := config.DB.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []model.Book
	for rows.Next() {
		var b model.Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.Price, &b.Sales, &b.Stock, &b.ImgPath, &b.CategoryID, &b.CategoryName)
		books = append(books, b)
	}
	return books, nil
}

// GetCategoryList 获取分类列表
func GetCategoryList() ([]map[string]interface{}, error) {
	query := "SELECT id,name FROM categories"
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		list = append(list, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}
	return list, nil
}

// 查询所有图书 分批处理，防止一次性把所有书籍都加载到内存或一次性往 Redis 写太多
func GetBooksByOffsetLimit(offset, limit int) ([]model.Book, error) {

	sqlStr := "SELECT id, stock FROM books ORDER BY id ASC LIMIT ? OFFSET ?"

	rows, err := config.DB.Query(sqlStr, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(&b.ID, &b.Stock); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}
