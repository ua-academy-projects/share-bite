package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sony/gobreaker"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/pkg/closer"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/database/txmanager"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	outboxpkg "github.com/ua-academy-projects/share-bite/pkg/outbox"
	"github.com/ua-academy-projects/share-bite/pkg/resilience"
)

func main() {
	baseCtx := context.Background()
	ctx, stop := signal.NotifyContext(baseCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// for local development only; .env is optional
	if err := config.Load(".env"); err != nil && !os.IsNotExist(err) {
		logger.Fatal(ctx, "config load:", err)
	}

	client, err := pg.NewClient(ctx, config.Config().Postgres.Dsn())
	if err != nil {
		logger.Fatal(ctx, "new database client:", err)
	}
	if err := client.DB().Ping(ctx); err != nil {
		logger.Fatal(ctx, "ping database:", err)
	}
	closer.Add(func(ctx context.Context) error {
		client.Close()
		return nil
	})
	closer.SetShutdownTimeout(5 * time.Second)

	topicArn := os.Getenv("OUTBOX_SNS_TOPIC_ARN")
	if topicArn == "" {
		logger.Fatal(ctx, "OUTBOX_SNS_TOPIC_ARN is required")
	}

	snsPub, err := outboxpkg.NewSNSPublisher(ctx, topicArn)
	if err != nil {
		logger.Fatal(ctx, "new sns publisher:", err)
	}

	store := outboxpkg.NewStore(client.DB())
	relay := outboxpkg.NewRelay(
		txmanager.NewTransactionManager(client.DB()),
		store,
		snsPub,
		outboxpkg.WithRelayPolicy(resilience.Policy{
			RetryConfig: resilience.RetryConfig{
				InitialInterval:     25 * time.Millisecond,
				RandomizationFactor: 0.2,
				Multiplier:          2,
				MaxInterval:         500 * time.Millisecond,
				MaxElapsedTime:      3 * time.Second,
			},
			Breaker: resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
				Name:        "outbox-sns-publish",
				MaxRequests: 1,
				Interval:    10 * time.Second,
				Timeout:     5 * time.Second,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.ConsecutiveFailures >= 10
				},
			}),
		}),
		outboxpkg.WithRelayPollInterval(2*time.Second),
	)

	go func() {
		if err := relay.Run(ctx); err != nil && err != context.Canceled {
			logger.Error(ctx, "outbox relay stopped:", err)
			os.Exit(1)
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if n, err := store.CleanupStuckProcessing(ctx, 5*time.Minute); err != nil {
					logger.ErrorKV(ctx, "failed to cleanup stuck outbox processing rows", "error", err)
				} else if n > 0 {
					logger.InfoKV(ctx, "cleaned up stuck outbox processing rows", "count", n)
				}
			}
		}
	}()

	closer.Wait()
}
