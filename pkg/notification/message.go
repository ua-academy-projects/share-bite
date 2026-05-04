package notification

import (
	"crypto/sha1"
	"fmt"
	"time"
)

type EventType string

const (
	PostLiked          EventType = "post_liked"
	InvitationReceived EventType = "invitation_received"
)

type Message struct {
	EventID     string    `json:"event_id"`
	EventType   EventType `json:"event_type"`
	RecipientID string    `json:"recipient_id"`
	ActorID     string    `json:"actor_id"`
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewEventID(parts ...string) string {
	h := sha1.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

func NewMessage(eventType EventType, recipientID, actorID, entityType, entityID string, createdAt time.Time) Message {
	return Message{
		EventID:     NewEventID(string(eventType), recipientID, actorID, entityType, entityID),
		EventType:   eventType,
		RecipientID: recipientID,
		ActorID:     actorID,
		EntityType:  entityType,
		EntityID:    entityID,
		CreatedAt:   createdAt,
	}
}
