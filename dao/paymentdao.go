package dao

import (
	"bookstore/model"
	"database/sql"
)

func InsertPaymentTx(tx *sql.Tx, p *model.Payment) error {
	sqlStr := `
		INSERT INTO payments(order_no,order_id,pay_method,pay_time,amount,status)
		VALUES(?, ?, ?, ?, ?, ?)
	`
	_, err := tx.Exec(sqlStr,
		p.OrderNo,
		p.OrderID,
		p.PayMethod,
		p.PayTime,
		p.Amount,
		p.Status,
	)
	return err
}
