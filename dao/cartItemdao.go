package dao

//
//// AddCartItem 向购物项表中插入购物项
//func AddCartItem(cartItem *model.CartItem) error {
//	sqlStr := "INSERT INTO cart_items(count,amount,book_id,cart_id) VALUES (?,?,?,?)"
//	_, err := config.DB.Exec(sqlStr, cartItem.Count, cartItem.GetAmount(), cartItem.Book.ID, cartItem.CartID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// GetCartItemByBookIDAndCartID 根据图书的id和购物车的id获取对应的购物项
//func GetCartItemByBookIDAndCartID(bookID string, cartID string) (*model.CartItem, error) {
//	sqlStr := "SELECT id,count,amount,cart_id FROM cart_items WHERE book_id = ? AND cart_id = ?"
//	row := config.DB.QueryRow(sqlStr, bookID, cartID)
//	cartItem := &model.CartItem{}
//
//	err := row.Scan(&cartItem.ID, &cartItem.Count, &cartItem.Amount, &cartItem.CartID)
//	if err != nil {
//		return nil, err
//	}
//
//	book, _ := GetBookByID(bookID)
//	cartItem.Book = book
//	return cartItem, nil
//
//}
//
//// UpdateBookCount 根据购物项中的相关信息更新购物项中图书的数量和金额小计
//func UpdateBookCount(cartItem *model.CartItem) error {
//	sql := "UPDATE cart_items SET count = ?, amount = ? WHERE book_id = ? AND cart_id = ?"
//	_, err := config.DB.Exec(sql, cartItem.Count, cartItem.GetAmount(), cartItem.Book.ID, cartItem.CartID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// GetCartItemsByCartID 根据购物车的id获取购物车中所有的购物项
//func GetCartItemsByCartID(cartID string) ([]*model.CartItem, error) {
//	sqlStr := "SELECT id,count,amount,book_id,cart_id FROM cart_items WHERE cart_id = ? "
//	rows, err := config.DB.Query(sqlStr, cartID)
//	if err != nil {
//		return nil, err
//	}
//	var cartItems []*model.CartItem
//	for rows.Next() {
//		var bookID string
//		cartItem := &model.CartItem{}
//		err2 := rows.Scan(&cartItem.ID, &cartItem.Count, &cartItem.Amount, &bookID, &cartItem.CartID)
//		if err2 != nil {
//			return nil, err2
//		}
//		book, _ := GetBookByID(bookID)
//		cartItem.Book = book
//		cartItems = append(cartItems, cartItem)
//	}
//	return cartItems, nil
//}
//
//// DeleteCartItemsByCartID 根据购物车的id删除所有的购物项
//func DeleteCartItemsByCartID(cartID string) error {
//	sql := "DELETE FROM cart_items  WHERE  cart_id = ?"
//	_, err := config.DB.Exec(sql, cartID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// DeleteCartItemByID 根据购物项的id删除购物项
//func DeleteCartItemByID(cartItemID string) error {
//	sql := "DELETE FROM cart_items WHERE id = ?"
//	_, err := config.DB.Exec(sql, cartItemID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
