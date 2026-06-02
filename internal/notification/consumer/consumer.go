package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type Processor interface {
	ProcessMessage(ctx context.Context, msg notification.Message) error
}

type SQSAPI interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

type Consumer struct {
	client      SQSAPI
	queueURL    string
	processor   Processor
	batchSize   int32
	waitSeconds int32
	pause       time.Duration
}

func New(client SQSAPI, queueURL string, processor Processor) *Consumer {
	return &Consumer{
		client:      client,
		queueURL:    queueURL,
		processor:   processor,
		batchSize:   10,
		waitSeconds: 20,
		pause:       2 * time.Second,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("notification consumer is nil")
	}
	if c.client == nil || c.processor == nil || c.queueURL == "" {
		return fmt.Errorf("notification consumer is not configured")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		resp, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.queueURL),
			MaxNumberOfMessages: c.batchSize,
			WaitTimeSeconds:     c.waitSeconds,
		})
		if err != nil {
			logger.ErrorKV(ctx, "receive notification sqs messages", "error", err)
			time.Sleep(c.pause)
			continue
		}

		for _, sqsMsg := range resp.Messages {
			if sqsMsg.Body == nil {
				continue
			}

			var msg notification.Message
			if err := json.Unmarshal([]byte(*sqsMsg.Body), &msg); err != nil {
				logger.ErrorKV(ctx, "unmarshal notification message", "message_id", aws.ToString(sqsMsg.MessageId), "error", err)
				continue
			}

			if err := c.processor.ProcessMessage(ctx, msg); err != nil {
				logger.ErrorKV(ctx, "process notification message", "message_id", aws.ToString(sqsMsg.MessageId), "error", err)
				continue
			}

			if sqsMsg.ReceiptHandle == nil {
				continue
			}

			if _, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(c.queueURL),
				ReceiptHandle: sqsMsg.ReceiptHandle,
			}); err != nil {
				logger.ErrorKV(ctx, "delete notification message", "message_id", aws.ToString(sqsMsg.MessageId), "error", err)
			}
		}
	}
}
