package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"
)

type Publisher interface {
	Publish(ctx context.Context, event Record) error
}

type Relay struct {
	txManager    database.TxManager
	store        Store
	publisher    Publisher
	policy       resilience.Policy
	batchSize    int
	pollInterval time.Duration
}

type RelayOption func(*Relay)

func WithRelayPolicy(policy resilience.Policy) RelayOption {
	return func(r *Relay) { r.policy = policy }
}

func WithRelayBatchSize(batchSize int) RelayOption {
	return func(r *Relay) { r.batchSize = batchSize }
}

func WithRelayPollInterval(interval time.Duration) RelayOption {
	return func(r *Relay) { r.pollInterval = interval }
}

func NewRelay(txManager database.TxManager, store Store, publisher Publisher, opts ...RelayOption) *Relay {
	r := &Relay{
		txManager:    txManager,
		store:        store,
		publisher:    publisher,
		batchSize:    100,
		pollInterval: 2 * time.Second,
		policy: resilience.Policy{
			RetryConfig: resilience.RetryConfig{
				InitialInterval:     50 * time.Millisecond,
				RandomizationFactor: 0.2,
				Multiplier:          2,
				MaxInterval:         1 * time.Second,
				MaxElapsedTime:      5 * time.Second,
			},
			Breaker: resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
				Name:        "outbox-relay",
				MaxRequests: 1,
				Interval:    10 * time.Second,
				Timeout:     5 * time.Second,
				ReadyToTrip: func(counts gobreaker.Counts) bool { return counts.ConsecutiveFailures >= 10 },
			}),
		},
	}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

func (r *Relay) Run(ctx context.Context) error {
	if r == nil {
		return fmt.Errorf("outbox relay is nil")
	}
	logger.InfoKV(ctx, "outbox relay started", "poll_interval", r.pollInterval, "batch_size", r.batchSize)
	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	for {
		if err := r.ProcessOnce(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.ErrorKV(ctx, "outbox relay cycle failed", "error", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (r *Relay) ProcessOnce(ctx context.Context) error {
	if r == nil {
		return fmt.Errorf("outbox relay is nil")
	}

	records, err := r.store.FetchPending(ctx, r.batchSize)
	if err != nil {
		return err
	}
	if len(records) > 0 {
		logger.InfoKV(ctx, "outbox relay fetched records", "count", len(records))
	}
	for _, record := range records {
		rec := record
		err := r.policy.Execute(ctx, func() error {
			return r.publisher.Publish(ctx, rec)
		})
		if err != nil {
			logger.ErrorKV(ctx, "failed to publish outbox record after retries",
				"record_id", rec.ID,
				"error", err)
			continue
		}

		if err := r.store.MarkProcessed(ctx, rec.ID); err != nil {
			return fmt.Errorf("mark processed %s: %w", rec.ID, err)
		}

		logger.InfoKV(ctx,
			"outbox record processed",
			"record_id", rec.ID,
			"event_id", rec.Payload.EventID,
			"event_type", rec.EventType,
		)
	}

	return nil
}
