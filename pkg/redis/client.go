package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr string, password string, db int, tlsEnabled bool) (*goredis.Client, error) {
	opts := &goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,

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
	}

	if tlsEnabled {
		opts.TLSConfig = &tls.Config{}
	}

	rdb := goredis.NewClient(opts)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return rdb, nil

}
