package model

// Cart 购物车结构体
type Cart struct {
	ID          int         `json:"id"`
	TotalCount  int         `json:"total_count"`
	TotalAmount float64     `json:"total_amount"`
	UserID      int         `json:"user_id"`
	CartItems   []*CartItem `json:"items,omitempty"`
}

// GetTotalCount 获取购物车中图书的总数量
func (c *Cart) GetTotalCount() int {
	var totalCount int
	for _, v := range c.CartItems {
		totalCount = totalCount + v.Count
	}
	return totalCount
}

// GetTotalAmount 获取购物车中图书的总金额
func (c *Cart) GetTotalAmount() float64 {
	var totalAmount float64
	for _, v := range c.CartItems {
		totalAmount = totalAmount + v.GetAmount()
	}
	return totalAmount
}
