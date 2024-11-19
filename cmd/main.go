package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MatthewAraujo/vacation-backend/cmd/api"
	configs "github.com/MatthewAraujo/vacation-backend/config"
	database "github.com/MatthewAraujo/vacation-backend/db"
	"github.com/MatthewAraujo/vacation-backend/pkg/assert"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	redisCfg := redis.Options{
		DB:       configs.Envs.Redis.Database,
		Password: configs.Envs.Redis.Password,
		Addr:     fmt.Sprintf("%s:%s", configs.Envs.Redis.Address, configs.Envs.Redis.Port),
	}

	db, err := database.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	redis := database.NewRedisStorage(redisCfg)

	err = healthRedis(redis)
	assert.NoError(err, "Redis is offline")
	log.Printf("Connect to redis")

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db, redis)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}

func healthRedis(redisClient *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
