package outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSPublisher struct {
	client   *sns.Client
	topicArn string
}

func NewSNSPublisher(ctx context.Context, topicArn string) (*SNSPublisher, error) {
	if topicArn == "" {
		return nil, fmt.Errorf("sns topic arn is required")
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	cli := sns.NewFromConfig(cfg)
	return &SNSPublisher{client: cli, topicArn: topicArn}, nil
}

func (p *SNSPublisher) Publish(ctx context.Context, event Message) error {
	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal sns message: %w", err)
	}

	_, err = p.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(p.topicArn),
		Message:  aws.String(string(b)),
	})
	if err != nil {
		return fmt.Errorf("sns publish: %w", err)
	}
	return nil
}
