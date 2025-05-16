package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func New(host, port, password string, db int) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	return redisClient, redisClient.Ping(context.Background()).Err()
}
