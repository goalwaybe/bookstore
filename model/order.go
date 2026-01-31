package model

import "time"

// Order 结构
type Order struct {
	ID          int          `json:"id"`
	CreateTime  time.Time    `json:"create_time"`
	TotalCount  int          `json:"total_count"`
	TotalAmount float64      `json:"total_amount"`
	State       int8         `json:"state"`
	UserID      int          `json:"user_id"`
	AddressID   int          `json:"address_id"`
	OrderItem   []*OrderItem `json:"order_item"`
	OrderNo     string       `json:"order_no"`
	Remark      string       `json:"remark"`
}
