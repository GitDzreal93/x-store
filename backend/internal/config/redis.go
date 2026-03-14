package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	cfg := Global.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := RDB.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("[Redis] Failed to connect: %v", err)
	}
	log.Printf("[Redis] Connected to %s", cfg.Addr())
}
