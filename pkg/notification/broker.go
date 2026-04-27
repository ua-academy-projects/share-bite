package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	redisclient "github.com/ua-academy-projects/share-bite/pkg/redis"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"
)

type client struct {
	rdb           *goredis.Client
	publishPolicy resilience.Policy
}

type Option func(*client)

type Publisher interface {
	Publish(ctx context.Context, ch string, msg Message) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, ch string) (<-chan Message, error)
}

func WithPublishPolicy(policy resilience.Policy) Option {
	return func(c *client) {
		c.publishPolicy = policy
	}
}

func NewBroker(rdb *goredis.Client, opts ...Option) *client {
	out := &client{
		rdb: rdb,
		publishPolicy: resilience.Policy{
			RetryConfig: resilience.RetryConfig{
				InitialInterval:     50 * time.Millisecond,
				RandomizationFactor: 0.2,
				Multiplier:          2,
				MaxInterval:         1 * time.Second,
				MaxElapsedTime:      5 * time.Second,
			},
		},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(out)
		}
	}

	return out
}

func (c *client) Publish(ctx context.Context, ch string, msg Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	operation := func() error {
		err := c.rdb.Publish(ctx, ch, msgBytes).Err()
		if redisclient.IsPermanentRedisError(err) {
			return resilience.Permanent(err)
		}

		return err
	}

	policy := c.publishPolicy
	configuredNotify := policy.RetryNotify
	policy.RetryNotify = func(err error, delay time.Duration) {
		if configuredNotify != nil {
			configuredNotify(err, delay)
			return
		}

		logger.Warnf(ctx, "Redis publish failed, retrying in %v: %v", delay, err)
	}

	err = policy.Execute(ctx, operation)
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
