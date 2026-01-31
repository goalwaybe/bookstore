package config

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func InitSnowflake(nodeID int64) {
	var err error
	Node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("snowflake 初始化失败: %v", err)
	}
}
