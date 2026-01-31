package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	cfg := Conf.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ 打开数据库失败：%v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)                   // 适当增加空闲连接
	db.SetConnMaxLifetime(110 * time.Second) // 缩短最大生命周期
	db.SetConnMaxIdleTime(30 * time.Second)  // 缩短空闲超时

	DB = db
	log.Println("✅ 数据库初始化成功")
}
