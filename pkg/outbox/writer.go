package outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ua-academy-projects/share-bite/pkg/database"
)

const DefaultSourceService = "guest-api"

// Event represents a row to be inserted into outbox table.
type Event struct {
	EventType     string
	Payload       any
	SourceService string
}

type Writer interface {
	Enqueue(ctx context.Context, event Event) error
}

type SQLWriter struct {
	db database.QueryExecer
}

func NewWriter(db database.QueryExecer) *SQLWriter {
	return &SQLWriter{db: db}
}

func (w *SQLWriter) Enqueue(ctx context.Context, event Event) error {
	if w == nil {
		return fmt.Errorf("outbox writer is nil")
	}
	if event.EventType == "" {
		return fmt.Errorf("outbox event type is required")
	}
	if event.SourceService == "" {
		return fmt.Errorf("outbox source service is required")
	}

	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("marshal outbox payload: %w", err)
	}

	q := database.Query{
		Name: "outbox_writer.Enqueue",
		Sql: `
			INSERT INTO outbox (event_type, payload, source_service)
			VALUES ($1, $2::jsonb, $3)
		`,
	}

	if _, err := w.db.ExecContext(ctx, q, event.EventType, string(payloadBytes), event.SourceService); err != nil {
		return fmt.Errorf("insert outbox event: %w", err)
	}

	return nil
}
