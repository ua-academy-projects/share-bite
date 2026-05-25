package outbox

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
)

const (
	EventTypePostLiked              = "post_liked"
	EventTypePostPublished          = "post_published"
	EventTypeRegistrationConfirmed  = "registration_confirmed"
	EventTypePasswordResetRequested = "password_reset_requested"
)

type Message struct {
	EventID     string         `json:"event_id"`
	EventType   string         `json:"event_type"`
	RecipientID string         `json:"recipient_id"`
	ActorID     string         `json:"actor_id"`
	EntityType  string         `json:"entity_type"`
	EntityID    string         `json:"entity_id"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

func NewEventID(parts ...string) string {
	h := sha1.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}
