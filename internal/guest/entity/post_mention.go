package entity

import "time"

type PostMention struct {
	PostID     string
	CustomerID string
	CreatedAt  time.Time
}
