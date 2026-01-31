package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type AppConfig struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

type JWTConfig struct {
	SecretKey     string `yaml:"secret_key"`
	ExpireSeconds int    `yaml:"expire_seconds"`
}

type Config struct {
	App    AppConfig    `yaml:"app"`
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	JWT    JWTConfig    `yaml:"jwt"`
}

var Conf Config

func InitConfig() {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("❌ 无法读取配置文件：%v", err)
	}

	if err := yaml.Unmarshal(data, &Conf); err != nil {
		log.Fatalf(" ❌ 无法解析配置文件: %v", err)
	}

	log.Println("✅ 配置文件加载成功")
}
