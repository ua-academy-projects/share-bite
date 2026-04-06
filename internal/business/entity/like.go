package entity

import "time"

type Like struct {
	ID         int64
	PostID     int64
	CustomerID string
	CreatedAt  time.Time
}
