package db

import (
	"github.com/MatthewAraujo/vacation-backend/pkg/assert"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func NewRedisStorage(cfg redis.Options) *redis.Client {
	assert.NotNil(cfg.Addr, "Redis address cannot be nil")
	assert.NotNil(cfg.Password, "Redis password cannot be nil")
	assert.NotNil(cfg.DB, "Redis DB cannot be nil")

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return rdb
}
