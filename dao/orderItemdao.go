package dao

import (
	"bookstore/config"
	"bookstore/model"
	"errors"
)

// CreateOrderItem
func CreateOrderItem(orderID, bookID, count int, price float64, orderNo string) error {

	book, err := GetBookDetail(bookID)
	if err != nil || book == nil {
		return errors.New("书籍不存在")
	}

	amount := float64(count) * price
	sqlStr := `
		INSERT INTO order_items (count,amount,title,author,price,img_path,order_id,book_id,order_no)
		VALUES (?,?,?,?,?,?,?,?,?)
	`
	_, err = config.DB.Exec(sqlStr,
		count,
		amount,
		book.Title,
		book.Author,
		book.Price,
		book.ImgPath,
		orderID,
		bookID,
		orderNo,
	)
	return err
}

// GetOrderItemsByOrderID 查询订单明细
func GetOrderItemsByOrderID(orderID int) ([]*model.OrderItem, error) {
	sqlStr := `
		SELECT id,book_id,count,amount,title,author,price,img_path,order_id
		FROM order_items
		WHERE order_id = ?
     `
	rows, err := config.DB.Query(sqlStr, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*model.OrderItem
	for rows.Next() {
		item := &model.OrderItem{}
		err := rows.Scan(
			&item.ID,
			&item.BookID,
			&item.Count,
			&item.Amount,
			&item.Title,
			&item.Author,
			&item.Price,
			&item.ImgPath,
			&item.OrderID,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
