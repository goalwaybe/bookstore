package cron

import (
	"bookstore/service"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCronJobs() {
	c := cron.New()

	stockService := &service.StockSyncService{}

	_, err := c.AddFunc("@every 10m", func() {
		fmt.Println(">>> [Cron] 开始全量同步库存...")
		if err := stockService.SyncAllStock(); err != nil {
			fmt.Println("[Cron] 库存同步失败:", err)
		} else {
			fmt.Println("[Cron] 库存同步成功")
		}
	})
	if err != nil {
		return
	}

	c.Start()
}
