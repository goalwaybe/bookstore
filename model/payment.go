package model

import "time"

type Payment struct {
	ID        int       `json:"id"`
	OrderNo   string    `json:"order_no"`
	OrderID   int       `json:"order_id"`
	PayMethod string    `json:"pay_method"`
	PayTime   time.Time `json:"pay_time"`
	Amount    float64   `json:"amount"`
	Status    int8      `json:"status"`
}
