package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ua-academy-projects/share-bite/internal/notification/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification entity.Notification) (bool, error)
	GetHistory(ctx context.Context, recipientID string, limit, offset int) ([]entity.Notification, error)
	MarkAsRead(ctx context.Context, recipientID string, notificationIDs []string) error
}

type SQLRepository struct {
	db database.DB
}

func New(db database.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Save(ctx context.Context, notification entity.Notification) (bool, error) {
	metadataJSON, err := json.Marshal(notification.Metadata)
	if err != nil {
		return false, fmt.Errorf("marshal metadata: %w", err)
	}

	q := database.Query{
		Name: "notification_repository.Save",
		Sql: `
			INSERT INTO notifications_history (
				notification_id,
				recipient_id,
				event_type,
				entity_id,
				metadata,
				created_at
			)
			VALUES ($1, $2, $3, $4, $5::jsonb, $6)
			ON CONFLICT (notification_id) DO NOTHING
		`,
	}

	tag, err := r.db.ExecContext(ctx, q, notification.NotificationID, notification.RecipientID, notification.EventType, notification.EntityID, string(metadataJSON), notification.CreatedAt)
	if err != nil {
		return false, fmt.Errorf("save notification: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}

func (r *SQLRepository) GetHistory(ctx context.Context, recipientID string, limit, offset int) ([]entity.Notification, error) {

	q := database.Query{
		Name: "notification_repository.GetHistory",
		Sql: `
			SELECT id, notification_id, recipient_id, event_type, entity_id, metadata, is_read, created_at, read_at
			FROM notifications_history
			WHERE recipient_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2 OFFSET $3
		`,
	}

	rows, err := r.db.QueryContext(ctx, q, recipientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get notification history: %w", err)
	}
	defer rows.Close()

	var items []entity.Notification
	if err := pgxscan.ScanAll(&items, rows); err != nil {
		return nil, fmt.Errorf("scan notification history: %w", err)
	}

	return items, nil
}

func (r *SQLRepository) MarkAsRead(ctx context.Context, recipientID string, notificationIDs []string) error {
	q := database.Query{
		Name: "notification_repository.MarkAsRead",
		Sql: `
			UPDATE notifications_history
			SET is_read = TRUE,
			    read_at = NOW()
			WHERE recipient_id = $1
			  AND notification_id = ANY($2)
		`,
	}

	if _, err := r.db.ExecContext(ctx, q, recipientID, notificationIDs); err != nil {
		return fmt.Errorf("mark notification as read: %w", err)
	}

	return nil
}

var _ NotificationRepository = (*SQLRepository)(nil)
