package outbox

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
)

const (
	EventTypePostLiked              = "post_liked"
	EventTypePostCommented          = "post_commented"
	EventTypePostMentioned          = "post_mentioned"
	EventTypePostInvitationReceived = "post_invitation_received"
	EventTypePostPublished          = "post_published"
	EventTypeRegistrationConfirmed  = "registration_confirmed"
)

type Message struct {
	EventID     string         `json:"eventID"`
	EventType   string         `json:"eventType"`
	RecipientID string         `json:"recipientID"`
	ActorID     string         `json:"actorID"`
	EntityType  string         `json:"entityType"`
	EntityID    string         `json:"entityID"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
}

func NewEventID(parts ...string) string {
	h := sha1.New()
	for _, p := range parts {
		h.Write([]byte(p))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}
