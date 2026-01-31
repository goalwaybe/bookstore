package dao

import (
	"bookstore/config"
	"bookstore/model"
	"database/sql"
	"errors"
)

// GetCartByUserID 根据用户的id从数据库中查询对应的购物车
func GetCartByUserID(userID int) (*model.Cart, error) {
	cart := &model.Cart{}
	sqlStr := "SELECT id,total_count,total_amount,user_id FROM carts WHERE user_id = ?"
	err := config.DB.QueryRow(sqlStr, userID).Scan(&cart.ID, &cart.TotalCount, &cart.TotalAmount, &cart.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`
		SELECT ci.id,ci.count,ci.amount,ci.cart_id,b.id,b.title,b.author,b.price
		FROM cart_items ci
		JOIN books b ON ci.book_id = b.id
		WHERE ci.cart_id = ?
`, cart.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.CartItem
		var book model.Book
		err := rows.Scan(&item.ID, &item.Count, &item.Amount, &item.CartID,
			&book.ID, &book.Title, &book.Author, &book.Price)
		if err != nil {
			return nil, err
		}
		item.Book = &book
		cart.CartItems = append(cart.CartItems, &item)
	}
	return cart, nil
}

// CreateCart 创建购物车
func CreateCart(userID int) (int64, error) {
	result, err := config.DB.Exec("INSERT INTO carts (user_id,total_count,total_amount) VALUES (?, 0, 0)", userID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetCartItem 获取购物车项
func GetCartItem(cartID, bookID int) (*model.CartItem, error) {
	item := &model.CartItem{
		Book: &model.Book{}, //关键！初始化Book
	}
	sqlStr := `
		SELECT ci.id,ci.count,ci.amount,ci.cart_id,b.id,b.title,b.author,b.price
		FROM cart_items ci
		JOIN books b ON ci.book_id = b.id
		WHERE ci.cart_id = ? AND ci.book_id = ?
`
	err := config.DB.QueryRow(sqlStr, cartID, bookID).
		Scan(&item.ID, &item.Count, &item.Amount, &item.CartID,
			&item.Book.ID, &item.Book.Title, &item.Book.Author, &item.Book.Price)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetCartItems 根据购物车ID获取所有购物项
func GetCartItems(cartID int) ([]*model.CartItem, error) {
	sqlStr := `
SELECT ci.id,ci.count,ci.amount,ci.book_id,ci.cart_id,
b.id,b.title,b.author,b.price
FROM cart_items ci
JOIN books b ON ci.book_id = b.id
WHERE ci.cart_id = ?
ORDER BY ci.id DESC
`
	rows, err := config.DB.Query(sqlStr, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*model.CartItem{}
	for rows.Next() {
		item := &model.CartItem{
			Book: &model.Book{},
		}
		err = rows.Scan(
			&item.ID,
			&item.Count,
			&item.Amount,
			&item.BookID,
			&item.CartID,
			&item.Book.ID,
			&item.Book.Title,
			&item.Book.Author,
			&item.Book.Price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// AddCartItem 添加购物项
func AddCartItem(cartID, bookID int, price float64, count int) error {
	amount := price * float64(count)
	_, err := config.DB.Exec("INSERT INTO cart_items (count,amount,book_id,cart_id,price) VALUES(?,?,?,?,?) ", count, amount, bookID, cartID, price)
	if err != nil {
		return err
	}
	_, err = config.DB.Exec("UPDATE carts SET total_count = total_count + ?, total_amount = total_amount + ? WHERE id = ? ", count, amount, cartID)
	return err
}

// UpdateCartItem 更新购物项数量
func UpdateCartItem(itemID, newCount int, price float64, cartID int) error {
	newAmount := price * float64(newCount)
	_, err := config.DB.Exec("UPDATE cart_items SET count = ?,amount = ?,price = ? WHERE id = ?", newCount, newAmount, price, itemID)
	if err != nil {
		return err
	}
	_, err = config.DB.Exec(`
		UPDATE carts
		SET total_count = (SELECT IFNULL(SUM(count),0) FROM cart_items WHERE cart_id = ?),
		    total_amount = (SELECT IFNULL(SUM(amount),0) FROM cart_items WHERE cart_id = ?)
		WHERE id = ?`, cartID, cartID, cartID)
	return err
}

// DeleteCartItem 删除购物项
func DeleteCartItem(cartID, bookID int) error {
	_, err := config.DB.Exec("DELETE FROM cart_items WHERE cart_id =? AND book_id = ?", cartID, bookID)
	if err != nil {
		return err
	}
	_, err = config.DB.Exec(`
		UPDATE carts
		SET total_count = (SELECT IFNULL(SUM(count),0) FROM cart_items WHERE cart_id = ?),
		total_amount = (SELECT IFNULL(SUM(amount),0) FROM cart_items WHERE cart_id = ?)
		WHERE id = ?`, cartID, cartID, cartID)
	return err
}

// ClearCart 清空购物车
func ClearCart(cartID int) error {
	_, err := config.DB.Exec("DELETE FROM cart_items WHERE cart_id = ? ", cartID)
	if err != nil {
		return err
	}
	_, err = config.DB.Exec("UPDATE carts SET total_count = 0,total_amount = 0 WHERE id = ?", cartID)
	return err
}
