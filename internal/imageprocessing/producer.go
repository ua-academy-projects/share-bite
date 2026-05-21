package imageprocessing

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Producer struct {
	client   *sqs.Client
	queueURL string
}

func NewProducer(
	client *sqs.Client,
	queueURL string,
) *Producer {
	return &Producer{
		client:   client,
		queueURL: queueURL,
	}
}

func (p *Producer) SendMessage(
	ctx context.Context,
	message ProcessImageMessage,
) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: aws.String(string(body)),
	})

	return err
}
