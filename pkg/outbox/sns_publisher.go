package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type SNSPublisher struct {
	client   *sns.Client
	topicArn string
}

func NewSNSPublisher(ctx context.Context, topicArn string) (*SNSPublisher, error) {
	if topicArn == "" {
		return nil, fmt.Errorf("sns topic arn is required")
	}
	topicRegion, err := parseTopicArn(topicArn)
	if err != nil {
		return nil, err
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	if !strings.EqualFold(cfg.Region, topicRegion) {
		logger.InfoKV(ctx,
			"sns client region adjusted to topic region",
			"config_region", cfg.Region,
			"topic_region", topicRegion,
		)
		cfg.Region = topicRegion
	}
	cli := sns.NewFromConfig(cfg)
	return &SNSPublisher{client: cli, topicArn: topicArn}, nil
}

func parseTopicArn(topicArn string) (region string, err error) {
	parts := strings.SplitN(topicArn, ":", 6)
	if len(parts) != 6 || parts[0] != "arn" || parts[2] != "sns" {
		return "", fmt.Errorf("sns topic arn must look like arn:aws:sns:<region>:<account>:<name>: %s", topicArn)
	}
	if parts[3] == "" {
		return "", fmt.Errorf("sns topic arn must include region: %s", topicArn)
	}
	return parts[3], nil
}

func (p *SNSPublisher) Publish(ctx context.Context, event Message) error {
	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal sns message: %w", err)
	}

	res, err := p.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(p.topicArn),
		Message:  aws.String(string(b)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.EventType),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("sns publish: %w", err)
	}

	logger.InfoKV(ctx,
		"sns publish succeeded",
		"topic_arn", p.topicArn,
		"message_id", aws.ToString(res.MessageId),
		"event_id", event.EventID,
		"event_type", event.EventType,
		"recipient_id", event.RecipientID,
	)

	return nil
}
