package main

import (
	"context"
	"encoding/json"
	postrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/post"
	"github.com/ua-academy-projects/share-bite/internal/storage"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ua-academy-projects/share-bite/internal/config"
	"github.com/ua-academy-projects/share-bite/internal/imageprocessing"
)

var processor *imageprocessing.Service

func init() {
	ctx := context.Background()

	err := config.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	storageClient, err := storage.NewStorageClient(ctx, config.Config().Storage)
	if err != nil {
		logger.Fatal(ctx, "init storage client:", err)
	}

	dbClient, err := pg.NewClient(
		ctx,
		config.Config().Postgres.Dsn(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	postRepository := postrepo.New(dbClient)

	processor = imageprocessing.NewService(
		storageClient,
		postRepository,
	)
}

func handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var msg imageprocessing.ProcessImageMessage

		err := json.Unmarshal(
			[]byte(record.Body),
			&msg,
		)
		if err != nil {
			return err
		}

		log.Printf(
			"processing image: image_id=%s s3_key=%s",
			msg.ImageID,
			msg.S3Key,
		)

		err = processor.ProcessImage(
			ctx,
			msg,
		)
		if err != nil {
			log.Printf(
				"failed to process image %s: %v",
				msg.ImageID,
				err,
			)

			continue
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
