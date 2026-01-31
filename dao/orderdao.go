package dao

import (
	"bookstore/config"
	"bookstore/model"
	"bookstore/utils"
	"database/sql"
	"time"
)

// CreateOrder 创建订单返回订单ID
func CreateOrder(userID, addressID int, totalAmount float64, totalCount int, remark string) (int, string, error) {
	orderNo := utils.GenerateOrderNo() // 生成雪花订单号

	sqlStr := `
		INSERT INTO orders (order_no,create_time,total_count,total_amount,state,user_id,address_id,remark)
		VALUES (?,?,?,?,0,?,?,?)
	`
	result, err := config.DB.Exec(sqlStr,
		orderNo,
		time.Now(),
		totalCount,
		totalAmount,
		userID,
		addressID,
		remark,
	)
	if err != nil {
		return 0, "", err
	}
	orderID64, err := result.LastInsertId()
	if err != nil {
		return 0, "", err
	}
	return int(orderID64), orderNo, nil
}

// GetOrdersByUserID 获取用户订单
func GetOrdersByUserID(userID, page, pageSize int) ([]*model.Order, error) {
	offset := (page - 1) * pageSize

	sqlStr := `
		SELECT id,create_time,total_count,total_amount,state,address_id
		FROM orders
		WHERE user_id = ?
		ORDER BY create_time DESC
		LIMIT ? OFFSET ?
`
	rows, err := config.DB.Query(sqlStr, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*model.Order{}

	for rows.Next() {
		order := &model.Order{}
		err := rows.Scan(
			&order.ID,
			&order.CreateTime,
			&order.TotalCount,
			&order.TotalAmount,
			&order.State,
			&order.AddressID,
		)
		if err != nil {
			return nil, err
		}

		items, err := GetOrderItemsByOrderID(order.ID)
		if err != nil {
			return nil, err
		}
		order.OrderItem = items

		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrderByID 根据用户订单ID获取订单详情
func GetOrderByID(orderID int) (*model.Order, error) {
	sqlStr := `
		SELECT id,create_time,total_count,total_amount,state,user_id,address_id,order_no
		FROM orders
		WHERE id = ?
`

	row := config.DB.QueryRow(sqlStr, orderID)
	order := &model.Order{}
	err := row.Scan(
		&order.ID,
		&order.CreateTime,
		&order.TotalCount,
		&order.TotalAmount,
		&order.State,
		&order.UserID,
		&order.AddressID,
		&order.OrderNo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	items, err := GetOrderItemsByOrderID(order.ID)
	if err != nil {
		return nil, err
	}
	order.OrderItem = items

	return order, nil
}

// CancelOrderByID 更新订单状态为已取消
func CancelOrderByID(orderID int) error {
	sqlStr := `
		UPDATE orders
		SET state = -1
		WHERE id = ? AND state IN (0,1)
`
	_, err := config.DB.Exec(sqlStr, orderID)
	return err
}

// PayOrderByID 更新订单状态为已支付
func PayOrderByIDTx(tx *sql.Tx, orderID int) error {
	sqlStr := `
		UPDATE orders
		SET state = 1
		WHERE id = ? AND state = 0
`
	_, err := tx.Exec(sqlStr, orderID)
	return err
}

// GetOrderByOrderNo 根据订单号查询到订单详情数据
func GetOrderByOrderNo(orderNo string) (*model.Order, error) {
	order := &model.Order{}
	sqlStr := "SELECT id,total_count, total_amount, state,user_id, address_id,remark,order_no FROM orders WHERE order_no = ? LIMIT 1"
	err := config.DB.QueryRow(sqlStr, orderNo).Scan(
		&order.ID,
		&order.TotalCount,
		&order.TotalAmount,
		&order.State,
		&order.UserID,
		&order.AddressID,
		&order.Remark,
		&order.OrderNo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有查到订单
		}
		return nil, err // 查询失败
	}
	return order, nil
}
