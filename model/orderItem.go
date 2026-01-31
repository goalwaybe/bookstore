package model

// OrderItem 结构
type OrderItem struct {
	ID      int     `json:"id"`
	Count   int     `json:"count"`
	Amount  float64 `json:"amount"`
	Title   string  `json:"title"`
	Author  string  `json:"author"`
	Price   float64 `json:"price"`
	ImgPath string  `json:"img_path"`
	OrderID string  `json:"order_id"`
	BookID  int     `json:"book_id"`
}
