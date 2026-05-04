package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type EventValidator interface {
	Validate(event notification.Message) error
}

type Processor interface {
	Process(ctx context.Context, event notification.Message) error
}

type Handler struct {
	validator EventValidator
	processor Processor
}

func New(validator EventValidator, processor Processor) *Handler {
	return &Handler{
		validator: validator,
		processor: processor,
	}
}

// HandleBatch processes an SQS event batch and returns failed message IDs for retry.
func (h *Handler) HandleBatch(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	response := events.SQSEventResponse{
		BatchItemFailures: []events.SQSBatchItemFailure{},
	}

	for _, record := range sqsEvent.Records {
		if err := h.handleRecord(ctx, record); err != nil {
			logger.ErrorKV(ctx, "failed to process SQS record", "message_id", record.MessageId, "error", err)
			response.BatchItemFailures = append(response.BatchItemFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}
	}

	return response, nil
}

func (h *Handler) handleRecord(ctx context.Context, record events.SQSMessage) error {
	var event notification.Message
	if err := json.Unmarshal([]byte(record.Body), &event); err != nil {
		return fmt.Errorf("unmarshal notification event: %w", err)
	}

	if err := h.validator.Validate(event); err != nil {
		return fmt.Errorf("validate notification event %q: %w", event.EventID, err)
	}

	if h.processor != nil {
		if err := h.processor.Process(ctx, event); err != nil {
			return fmt.Errorf("process notification event %q: %w", event.EventID, err)
		}
	}

	return nil
}
