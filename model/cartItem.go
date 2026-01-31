package model

// CartItem 购物项结构体
type CartItem struct {
	ID     int     `json:"id"`
	Price  float64 `json:"price"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
	BookID int     `json:"book_id"`
	CartID int     `json:"cart_id"`
	Book   *Book   `json:"book"`
}

// GetAmount 获取购物项中图书的金额小计，有图书的价格和图书的数量计算得到
func (ci *CartItem) GetAmount() float64 {
	price := ci.Book.Price
	return float64(ci.Count) * price
}
