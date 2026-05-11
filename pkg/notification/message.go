package notification

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	PostLiked             EventType = "post_liked"
	RegistrationConfirmed EventType = "registration_confirmed"
	InvitationReceived    EventType = "invitation_received"
)

type Message struct {
	EventID     string         `json:"eventID"`     // Unique event ID (SHA-256)
	EventType   EventType      `json:"eventType"`   // Event type (e.g. "post_liked")
	RecipientID string         `json:"recipientID"` // Who receives the notification (user ID)
	ActorID     string         `json:"actorID"`     // Who triggered the event (user ID)
	EntityType  string         `json:"entityType"`  // Type of entity (e.g. post)
	EntityID    string         `json:"entityID"`    // ID of the entity (e.g. post ID)
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"` // Event timestamp
}

func NewEventID() string {
	return uuid.NewString()
}

func NewMessage(eventType EventType, recipientID, actorID, entityType, entityID string, createdAt time.Time) Message {
	return NewMessageWithMetadata(eventType, recipientID, actorID, entityType, entityID, nil, createdAt)
}

func NewMessageWithMetadata(eventType EventType, recipientID, actorID, entityType, entityID string, metadata map[string]any, createdAt time.Time) Message {
	var clonedMetadata map[string]any
	if metadata != nil {
		clonedMetadata = make(map[string]any, len(metadata))
		for k, v := range metadata {
			clonedMetadata[k] = v
		}
	}

	return Message{
		EventID:     NewEventID(),
		EventType:   eventType,
		RecipientID: recipientID,
		ActorID:     actorID,
		EntityType:  entityType,
		EntityID:    entityID,
		Metadata:    clonedMetadata,
		CreatedAt:   createdAt,
	}
}
