package rediskit

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type Redisconfig struct {
	Addr     string
	Password string
	Database int
}

type RedisClient struct {
	*redis.Client
	closeFunc func()
}

func (c *RedisClient) Close() error {
	if c.closeFunc != nil {
		c.closeFunc()
	}

	return c.Client.Close()
}

func NewRedisClient(ctx context.Context, conf *Redisconfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.Database,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	return &RedisClient{
		Client: client,
	}
}
