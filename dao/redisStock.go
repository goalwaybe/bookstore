package dao

import (
	"bookstore/config"
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

var ctx = context.Background()

// Redis Key
const (
	RealStockKey   = "book_stock_real"   //真实库存
	FrozenStockKey = "book_stock_frozen" //预扣库存
)

// -----------------------------------------------
// 初始化库存
// 后台上架或修改库存时调用
// -----------------------------------------------
// InitStock 初始化库存: 初始化真实库存
func InitStock(bookID int, stock int) error {
	field := strconv.Itoa(bookID)
	return config.RedisClient.HSet(ctx, RealStockKey, field, stock).Err()
}

// UpdateRealStock 修改真实库存(后台修改库存)
func UpdateRealStock(bookID int, stock int) error {
	field := strconv.Itoa(bookID)
	return config.RedisClient.HSet(ctx, RealStockKey, field, stock).Err()
}

// DeleteRealStock 删除图书时删除库存记录
func DeleteRealStock(bookID int) error {
	field := strconv.Itoa(bookID)
	pipe := config.RedisClient.TxPipeline()
	pipe.HDel(ctx, RealStockKey, field)
	pipe.HDel(ctx, FrozenStockKey, field)
	_, err := pipe.Exec(ctx)
	return err
}

// -----------------------------------------------
// 用户下单用方法: 预扣库存
// PreDeductStock 预扣库存
// -----------------------------------------------
func PreDeductStock(bookID int, count int) error {

	field := strconv.Itoa(bookID)

	lua := `
	local real = tonumber(redis.call("HGET",KEYS[1],ARGV[1]) or "-1")
	if real == -1 then
		return -1
	end
	local frozen = tonumber(redis.call("HGET",KEYS[2],ARGV[1]) or "0")
	local available = real - frozen
	if available < tonumber(ARGV[2]) then
		return 0
	end
	redis.call("HINCRBY",KEYS[2],ARGV[1],ARGV[2])
	return 1
	`
	res, err := config.RedisClient.Eval(ctx, lua, []string{RealStockKey, FrozenStockKey}, field, count).Int()
	if err != nil {
		return err
	}
	if res == -1 {
		return fmt.Errorf("商品不存在")
	}
	if res == 0 {
		return fmt.Errorf("库存不足")
	}

	return nil
}

// RollbackStock 取消订单:归还预扣库存
func RollbackStock(bookID int, count int) error {
	field := strconv.Itoa(bookID)
	return config.RedisClient.HIncrBy(ctx, FrozenStockKey, field, int64(-count)).Err()
}

// ConfirmStock 确认库存(支付后写入 MySQL)
func ConfirmStockTx(tx *sql.Tx, bookID int, count int) error {
	// UPDATE books SET stock = stock - count WHERE id = ? AND stock >= count

	field := strconv.Itoa(bookID)

	lua := `
	local frozen = tonumber(redis.call("HGET",KEYS[2],ARGV[1]) or "0")
	if frozen < tonumber(ARGV[2]) then
		return 0
	end
	redis.call("HINCRBY",KEYS[2],ARGV[1],-tonumber(ARGV[2]))
	redis.call("HINCRBY",KEYS[1],ARGV[1],-tonumber(ARGV[2]))
	return 1
	`
	res, err := config.RedisClient.Eval(ctx, lua, []string{RealStockKey, FrozenStockKey}, field, count).Int()
	if err != nil {
		return err
	}
	if res == 0 {
		return fmt.Errorf("冻结库存不足")
	}

	// -----------------------------------------------
	// Mysql 扣减数据库库存
	// -----------------------------------------------

	sqlStr := `
		UPDATE books
		SET stock = stock - ?
		WHERE id = ? AND stock >= ?
	`
	result, err := tx.Exec(sqlStr, count, bookID, count)
	if err != nil {
		config.RedisClient.HIncrBy(ctx, RealStockKey, field, int64(count))
		config.RedisClient.HIncrBy(ctx, FrozenStockKey, field, int64(count))
		return fmt.Errorf("数据库库存扣减失败: %v", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		config.RedisClient.HIncrBy(ctx, RealStockKey, field, int64(count))
		config.RedisClient.HIncrBy(ctx, RealStockKey, field, int64(count))

		return fmt.Errorf("数据库库存不足")
	}

	return nil
}
