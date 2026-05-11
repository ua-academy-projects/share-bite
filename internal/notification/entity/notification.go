package entity

import (
	"time"

	notificationpkg "github.com/ua-academy-projects/share-bite/pkg/notification"
)

type Notification struct {
	ID             int64
	NotificationID string
	RecipientID    string
	EventType      string
	EntityID       string
	Metadata       map[string]any
	IsRead         bool
	CreatedAt      time.Time
	ReadAt         *time.Time
}

type NotificationDTO struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	EntityID  string         `json:"entityID"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

func (n Notification) ToDTO() NotificationDTO {
	return NotificationDTO{
		ID:        n.NotificationID,
		Type:      n.EventType,
		EntityID:  n.EntityID,
		Metadata:  n.Metadata,
		CreatedAt: n.CreatedAt,
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
