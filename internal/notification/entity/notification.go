package entity

import (
	"time"

	notificationpkg "github.com/ua-academy-projects/share-bite/pkg/notification"
)

type Notification struct {
	ID             int64          `db:"id"`
	NotificationID string         `db:"notification_id"`
	RecipientID    string         `db:"recipient_id"`
	EventType      string         `db:"event_type"`
	EntityID       string         `db:"entity_id"`
	Metadata       map[string]any `db:"metadata"`
	IsRead         bool           `db:"is_read"`
	CreatedAt      time.Time      `db:"created_at"`
	ReadAt         *time.Time     `db:"read_at"`
}

type NotificationDTO struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	EntityID  string         `json:"entityID"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	IsRead    bool           `json:"isRead"`
	CreatedAt time.Time      `json:"createdAt"`
	ReadAt    *time.Time     `json:"readAt,omitempty"`
}

func (n Notification) ToDTO() NotificationDTO {
	return NotificationDTO{
		ID:        n.NotificationID,
		Type:      n.EventType,
		EntityID:  n.EntityID,
		Metadata:  n.Metadata,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
		ReadAt:    n.ReadAt,
	}
}

func FromMessage(msg notificationpkg.Message) Notification {
	return Notification{
		NotificationID: msg.EventID,
		RecipientID:    msg.RecipientID,
		EventType:      string(msg.EventType),
		EntityID:       msg.EntityID,
		Metadata:       msg.Metadata,
		CreatedAt:      msg.CreatedAt,
	}
}
