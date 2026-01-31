package utils

import (
	"bookstore/config"
	"strconv"
	"time"
)

func GenerateOrderNo() string {
	if config.Node == nil {
		panic("Snowflake node 未初始化...")
	}

	// 获取当前时间的年月日时分秒(14位)
	timestamp := time.Now().Format("20060102150405")

	//获取雪花算法序列号(取后4位)
	sequence := config.Node.Generate().Int64() % 10000

	//格式 : ORD + 年月日时分秒(14位) + 序列号(4位) = 总长度21位
	return "ORD" + timestamp + strconv.FormatInt(sequence, 10)
}
