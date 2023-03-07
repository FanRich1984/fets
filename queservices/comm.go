package queservices

import "github.com/go-redis/redis/v8"

func NewTmpRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
}
