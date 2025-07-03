package database

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Params   string
	MaxConns int
}

// MySQLDSN генерирует строку подключения для MySQL
func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Params,
	)
}

// NewConfigFromEnv создает конфиг из переменных окружения
func NewConfigFromEnv() *Config {
	return &Config{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName:   os.Getenv("MYSQL_DATABASE"),
		Params:   "charset=utf8mb4&parseTime=True&loc=Local",
		MaxConns: 25,
	}
}