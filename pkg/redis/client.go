package redis

import (
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr string, password string) (*goredis.Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		PoolSize:        100,
		MinIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		PoolTimeout:     4 * time.Second,

		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})
	return rdb, nil

}
