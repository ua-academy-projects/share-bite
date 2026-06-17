package worker_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ua-academy-projects/share-bite/internal/notification/worker"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type mockValidator struct {
	validateFn func(event notification.Message) error
}

func (m *mockValidator) Validate(event notification.Message) error {
	if m.validateFn != nil {
		return m.validateFn(event)
	}
	return nil
}

type mockProcessor struct {
	processFn func(ctx context.Context, event notification.Message) error
}

func (m *mockProcessor) Process(ctx context.Context, event notification.Message) error {
	if m.processFn != nil {
		return m.processFn(ctx, event)
	}
	return nil
}

func TestHandler_HandleBatch_ValidEvent(t *testing.T) {
	t.Parallel()

	event := notification.Message{
		EventID:     "test-event-123",
		EventType:   notification.RegistrationConfirmed,
		RecipientID: "user-42",
		ActorID:     "user-99",
		EntityType:  "post",
		EntityID:    "post-123",
		CreatedAt:   time.Now(),
		Metadata: map[string]any{
			"email": "test@example.com",
		},
	}

	body, err := json.Marshal(event)
	require.NoError(t, err)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     "msg-001",
				ReceiptHandle: "handle-001",
				Body:          string(body),
			},
		},
	}

	processed := false
	validator := &mockValidator{}
	processor := &mockProcessor{
		processFn: func(ctx context.Context, e notification.Message) error {
			processed = true
			assert.Equal(t, event.EventID, e.EventID)
			return nil
		},
	}

	h := worker.New(validator, processor)
	response, err := h.HandleBatch(context.Background(), sqsEvent)

	require.NoError(t, err)
	assert.True(t, processed)
	assert.Empty(t, response.BatchItemFailures)
}

func TestHandler_HandleBatch_InvalidEvent(t *testing.T) {
	t.Parallel()

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     "msg-001",
				ReceiptHandle: "handle-001",
				Body:          `{"eventType":"registration_confirmed"}`,
			},
		},
	}

	validator := worker.NewDefaultValidator(notification.RegistrationConfirmed)
	processor := &mockProcessor{}

	h := worker.New(validator, processor)
	response, err := h.HandleBatch(context.Background(), sqsEvent)

	require.NoError(t, err)
	assert.Len(t, response.BatchItemFailures, 1)
	assert.Equal(t, "msg-001", response.BatchItemFailures[0].ItemIdentifier)
}

func TestHandler_HandleBatch_PartialFailure(t *testing.T) {
	t.Parallel()

	validEvent := notification.Message{
		EventID:     "test-event-valid",
		EventType:   notification.RegistrationConfirmed,
		RecipientID: "user-42",
		ActorID:     "user-99",
		EntityType:  "post",
		EntityID:    "post-123",
		CreatedAt:   time.Now(),
		Metadata: map[string]any{
			"email":    "test@example.com",
			"username": "test-user",
		},
	}

	validBody, err := json.Marshal(validEvent)
	require.NoError(t, err)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     "msg-valid",
				ReceiptHandle: "handle-valid",
				Body:          string(validBody),
			},
			{
				MessageId:     "msg-invalid",
				ReceiptHandle: "handle-invalid",
				Body:          `{"eventType":"unknown"}`,
			},
		},
	}

	processed := false
	validator := worker.NewDefaultValidator(notification.RegistrationConfirmed)
	processor := &mockProcessor{
		processFn: func(ctx context.Context, e notification.Message) error {
			processed = true
			return nil
		},
	}

	h := worker.New(validator, processor)
	response, err := h.HandleBatch(context.Background(), sqsEvent)

	require.NoError(t, err)
	assert.True(t, processed)
	assert.Len(t, response.BatchItemFailures, 1)
	assert.Equal(t, "msg-invalid", response.BatchItemFailures[0].ItemIdentifier)
}

func TestHandler_HandleBatch_MalformedJSON(t *testing.T) {
	t.Parallel()

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     "msg-malformed",
				ReceiptHandle: "handle-malformed",
				Body:          `{invalid json`,
			},
		},
	}

	validator := worker.NewDefaultValidator()
	processor := &mockProcessor{}

	h := worker.New(validator, processor)
	response, err := h.HandleBatch(context.Background(), sqsEvent)

	require.NoError(t, err)
	assert.Len(t, response.BatchItemFailures, 1)
	assert.Equal(t, "msg-malformed", response.BatchItemFailures[0].ItemIdentifier)
}

func TestHandler_HandleBatch_ProcessorError(t *testing.T) {
	t.Parallel()

	event := notification.Message{
		EventID:     "test-event-123",
		EventType:   notification.RegistrationConfirmed,
		RecipientID: "user-42",
		ActorID:     "user-99",
		EntityType:  "post",
		EntityID:    "post-123",
		CreatedAt:   time.Now(),
		Metadata: map[string]any{
			"email": "test@example.com",
		},
	}

	body, err := json.Marshal(event)
	require.NoError(t, err)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     "msg-001",
				ReceiptHandle: "handle-001",
				Body:          string(body),
			},
		},
	}

	validator := &mockValidator{}
	processor := &mockProcessor{
		processFn: func(ctx context.Context, e notification.Message) error {
			return assert.AnError
		},
	}

	h := worker.New(validator, processor)
	response, err := h.HandleBatch(context.Background(), sqsEvent)

	require.NoError(t, err)
	assert.Len(t, response.BatchItemFailures, 1)
	assert.Equal(t, "msg-001", response.BatchItemFailures[0].ItemIdentifier)
}

func TestHandler_HandleBatch_FromFixtures(t *testing.T) {
	t.Parallel()

	validBatchData, err := os.ReadFile("testdata/valid_batch.json")
	require.NoError(t, err)

	var validBatch events.SQSEvent
	err = json.Unmarshal(validBatchData, &validBatch)
	require.NoError(t, err)

	processed := false
	validator := worker.NewDefaultValidator(notification.RegistrationConfirmed)
	processor := &mockProcessor{
		processFn: func(ctx context.Context, e notification.Message) error {
			processed = true
			return nil
		},
	}

	h := worker.New(validator, processor)
	response, err := h.HandleBatch(context.Background(), validBatch)

	require.NoError(t, err)
	assert.True(t, processed)
	assert.Empty(t, response.BatchItemFailures)

	mixedBatchData, err := os.ReadFile("testdata/batch_with_invalid_event.json")
	require.NoError(t, err)

	var mixedBatch events.SQSEvent
	err = json.Unmarshal(mixedBatchData, &mixedBatch)
	require.NoError(t, err)

	processed = false
	response, err = h.HandleBatch(context.Background(), mixedBatch)

	require.NoError(t, err)
	assert.True(t, processed)
	assert.Len(t, response.BatchItemFailures, 1)
	assert.Equal(t, "msg-002", response.BatchItemFailures[0].ItemIdentifier)
}

func TestHandler_HandleBatch_PasswordResetEvent(t *testing.T) {
	t.Parallel()

	event := notification.Message{
		EventID:     "reset-event-123",
		EventType:   notification.PasswordResetRequested,
		RecipientID: "user-42",
		ActorID:     "user-42",
		EntityType:  "user",
		EntityID:    "user-42",
		CreatedAt:   time.Now(),
		Metadata: map[string]any{
			"email":       "test@example.com",
			"reset_token": "token-123",
		},
	}

	body, err := json.Marshal(event)
	require.NoError(t, err)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{{
			MessageId:     "msg-reset",
			ReceiptHandle: "handle-reset",
			Body:          string(body),
		}},
	}

	processed := false
	validator := worker.NewDefaultValidator(notification.RegistrationConfirmed, notification.PasswordResetRequested)
	processor := &mockProcessor{
		processFn: func(ctx context.Context, e notification.Message) error {
			processed = true
			assert.Equal(t, notification.PasswordResetRequested, e.EventType)
			return nil
		},
	}

	h := worker.New(validator, processor)
	response, err := h.HandleBatch(context.Background(), sqsEvent)

	require.NoError(t, err)
	assert.True(t, processed)
	assert.Empty(t, response.BatchItemFailures)
}
