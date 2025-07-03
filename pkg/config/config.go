package config

import (
	"fmt"
	"os"
	"strconv"
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

type RedisConfig struct {
    Address  string
    Password string
    DB       int
    // PoolSize int
}

func NewRedisConfigFromEnv() *RedisConfig {
    return &RedisConfig{
        Address: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
        DB:       getEnvAsInt("REDIS_DB", 0),
        // PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 10),
    }
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := os.Getenv(key)
    if valueStr == "" {
        return defaultValue
    }
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}