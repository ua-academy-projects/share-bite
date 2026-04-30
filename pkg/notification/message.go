package notification

import "time"

type EventType string

const (
	PostLiked EventType = "post_liked"
)

type Message struct {
	UserID    string    `json:"user_id"`
	Type      EventType `json:"type"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}
