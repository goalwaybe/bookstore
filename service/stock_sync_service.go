package service

import (
	"bookstore/config"
	"bookstore/dao"
	"bookstore/model"
	"context"
	"log"
	"strconv"
)

var ctx = context.Background()

type StockSyncService struct{}

var StockSync = &StockSyncService{}

// -----------------------------------------------
// 同步单个书籍库存
// -----------------------------------------------
func (s *StockSyncService) syncOneBookStock(book *model.Book) error {
	field := strconv.Itoa(int(book.ID))

	err := config.RedisClient.HSet(ctx, dao.RealStockKey, field, book.Stock).Err()
	if err != nil {
		log.Println("[StockSync] SyncOneBookStock Error")
		return err
	}
	return nil
}

// -----------------------------------------------
// 删除书籍库存(删除书籍时调用) real + frozen 都要删除
// -----------------------------------------------
func (s *StockSyncService) DeleteBookStock(bookID int64) error {
	field := strconv.FormatInt(bookID, 10)

	pipe := config.RedisClient.TxPipeline()
	pipe.HDel(ctx, dao.RealStockKey, field)
	pipe.HDel(ctx, dao.FrozenStockKey, field)

	_, err := pipe.Exec(ctx)

	if err != nil {
		log.Println("[StockSync] DeleteBookStock Error:", err)
		return err
	}
	return nil
}

// -----------------------------------------------
// 分批同步所有库存 (用于项目启动/定时任务)
// 只同步 real，不触碰 frozen
// -----------------------------------------------
func (s *StockSyncService) SyncAllStock() error {

	const batchSize = 100
	offset := 0

	for {
		//1.分批获取书籍
		books, err := dao.GetBooksByOffsetLimit(offset, batchSize)
		if err != nil {
			log.Println("[StockSync] DB Error : ", err)
			return err
		}

		//没有更多数据
		if len(books) == 0 {
			break
		}

		//2. Redis pipeline 提高批量性能
		pipe := config.RedisClient.Pipeline()

		for _, book := range books {
			field := strconv.Itoa(int(book.ID))
			pipe.HSet(ctx, dao.RealStockKey, field, book.Stock)
		}

		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Println("[StockSync] Pipeline Error:", err)
			return err
		}

		log.Printf("[StockSync] Synced batch: offset=%d, count=%d\n", offset, len(books))

		offset += batchSize
	}

	log.Println("[StockSync] SyncAllStock Finished")
	return nil
}
