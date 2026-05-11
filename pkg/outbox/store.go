package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Record struct {
	ID            string
	EventType     string
	Payload       Message
	SourceService string
	CreatedAt     time.Time
}

type Store interface {
	FetchPending(ctx context.Context, limit int) ([]Record, error)
	MarkProcessed(ctx context.Context, id string) error
	CleanupStuckProcessing(ctx context.Context, olderThan time.Duration) (int64, error)
}

type SQLStore struct {
	db database.DB
}

func NewStore(db database.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) FetchPending(ctx context.Context, limit int) ([]Record, error) {
	if s == nil {
		return nil, fmt.Errorf("outbox store is nil")
	}
	if limit <= 0 {
		limit = 100
	}
	q := database.Query{
		Name: "outbox_store.FetchPending",
		Sql: `
            UPDATE outbox
            SET status = 'processing', updated_at = NOW()
            WHERE id IN (
                SELECT id 
                FROM outbox
                WHERE status = 'pending'
                ORDER BY created_at ASC
                LIMIT $1
                FOR UPDATE SKIP LOCKED
            )
            RETURNING id, event_type, payload, source_service, created_at
        `,
	}

	rows, err := s.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch pending outbox rows: %w", err)
	}
	defer rows.Close()

	records := make([]Record, 0, limit)
	for rows.Next() {
		var rec Record
		var payloadBytes []byte
		if err := rows.Scan(&rec.ID, &rec.EventType, &payloadBytes, &rec.SourceService, &rec.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan pending outbox row: %w", err)
		}
		if err := json.Unmarshal(payloadBytes, &rec.Payload); err != nil {
			return nil, fmt.Errorf("unmarshal outbox payload %s: %w", rec.ID, err)
		}
		records = append(records, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending outbox rows: %w", err)
	}

	return records, nil
}

func (s *SQLStore) MarkProcessed(ctx context.Context, id string) error {
	if s == nil {
		return fmt.Errorf("outbox store is nil")
	}

	q := database.Query{
		Name: "outbox_store.MarkProcessed",
		Sql: `
			UPDATE outbox
			SET status = 'processed', processed_at = NOW(), updated_at = NOW()
			WHERE id = $1
		`,
	}

	if _, err := s.db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("mark outbox row processed: %w", err)
	}

	return nil
}

func (s *SQLStore) CleanupStuckProcessing(ctx context.Context, olderThan time.Duration) (int64, error) {
	if s == nil {
		return 0, fmt.Errorf("outbox store is nil")
	}
	cutoff := time.Now().Add(-olderThan)
	q := database.Query{
		Name: "outbox_store.CleanupStuckProcessing",
		Sql: `
			UPDATE outbox
			SET status = 'pending', updated_at = NOW()
			WHERE status = 'processing' AND updated_at < $1
		`,
	}
	res, err := s.db.ExecContext(ctx, q, cutoff)
	if err != nil {
		return 0, fmt.Errorf("cleanup stuck processing: %w", err)
	}
	return res.RowsAffected(), nil
}

var _ Store = (*SQLStore)(nil)
