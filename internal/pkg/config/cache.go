package config

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppConfig.REDIS_HOST, AppConfig.REDIS_PORT),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
