package model

// Address 收货地址模型
type Address struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	ReceiverName string `json:"receiver_name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	IsDefault    int8   `json:"is_default"`
}
