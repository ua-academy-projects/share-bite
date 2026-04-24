package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/redis/go-redis/v9"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type client struct {
	rdb *redis.Client
}

type Publisher interface {
	Publish(ctx context.Context, ch string, msg Message) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, ch string) (<-chan Message, error)
}

func NewBroker(rdb *redis.Client) *client {
	return &client{rdb: rdb}
}

func (c *client) Publish(ctx context.Context, ch string, msg Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	operation := func() error {
		err := c.rdb.Publish(ctx, ch, msgBytes).Err()
		if err != nil {
			return err
		}
		return nil
	}

	notify := func(err error, delay time.Duration) {
		logger.Warnf(ctx, "Redis publish failed, retrying in %v: %v", delay, err)
	}

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 50 * time.Millisecond
	b.MaxInterval = 1 * time.Second
	b.MaxElapsedTime = 5 * time.Second
	err = backoff.RetryNotify(
		operation,
		backoff.WithContext(b, ctx),
		notify,
	)

	if err != nil {
		return fmt.Errorf("redis publish with retry: %w", err)
	}

	return nil
}

func (c *client) Subscribe(ctx context.Context, ch string) (<-chan Message, error) {
	pubsub := c.rdb.Subscribe(ctx, ch)
	if _, err := pubsub.Receive(ctx); err != nil {
		return nil, fmt.Errorf("redis subscribe error: %v\n", err)
	}
	out := make(chan Message, 100)
	redisCh := pubsub.Channel()
	go func() {
		defer pubsub.Close()
		defer close(out)
		for msg := range redisCh {
			var payload Message
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				logger.WarnKV(ctx, "skip malformed notification payload", "payload", msg.Payload, "error", err)
				continue
			}
			out <- payload
		}
	}()
	return out, nil
}
