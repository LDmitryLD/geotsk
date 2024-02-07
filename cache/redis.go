package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host, port string) *redis.Client {
	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return client
}
