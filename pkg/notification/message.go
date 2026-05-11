package notification

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type EventType string

const (
	PostLiked             EventType = "post_liked"
	RegistrationConfirmed EventType = "registration_confirmed"
	InvitationReceived    EventType = "invitation_received"
)

type Message struct {
	EventID     string         `json:"event_id"`     // Unique event ID (SHA-256)
	EventType   EventType      `json:"event_type"`   // Event type (e.g. "post_liked")
	RecipientID string         `json:"recipient_id"` // Who receives the notification (user ID)
	ActorID     string         `json:"actor_id"`     // Who triggered the event (user ID)
	EntityType  string         `json:"entity_type"`  // Type of entity (e.g. post)
	EntityID    string         `json:"entity_id"`    // ID of the entity (e.g. post ID)
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"` // Event timestamp
}

func NewEventID(parts ...string) string {
	h := sha256.New()
	for i, p := range parts {
		if i > 0 {
			h.Write([]byte{0}) // null byte delimiter
		}
		h.Write([]byte(p))
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:32]
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
		EventID:     NewEventID(string(eventType), recipientID, actorID, entityType, entityID),
		EventType:   eventType,
		RecipientID: recipientID,
		ActorID:     actorID,
		EntityType:  entityType,
		EntityID:    entityID,
		Metadata:    clonedMetadata,
		CreatedAt:   createdAt,
	}
}
