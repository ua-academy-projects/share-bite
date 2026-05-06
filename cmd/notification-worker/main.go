package main

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/internal/notification/worker"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
	redispkg "github.com/ua-academy-projects/share-bite/pkg/redis"
)

func main() {
	ctx := context.Background()

	// load env/config (local .env or environment variables)
	if err := config.Load(".env"); err != nil {
		logger.Fatal(ctx, "config load:", err)
	}

	// initialize redis client
	rdb, err := redispkg.NewClient(
		config.Config().Redis.Addr(),
		config.Config().Redis.Password(),
		config.Config().Redis.DB(),
		config.Config().Redis.TLS(),
	)
	if err != nil {
		logger.Fatal(ctx, "new redis client:", err)
	}

	// producer (broker) used to publish notifications to subscribers
	broker := notification.NewBroker(rdb)

	validator := worker.NewDefaultValidator()
	processor := worker.NewPublisherProcessor(broker, 5*time.Second)
	h := worker.New(validator, processor)

	lambda.Start(h.HandleBatch)
}
