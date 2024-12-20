package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisInst *redis.Client

func InitRedis(address, password string, DB int) (RedisInst, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %v", err)
	}
	return client, nil
}
