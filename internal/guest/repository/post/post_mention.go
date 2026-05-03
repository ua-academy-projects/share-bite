package post

import "time"

type PostMention struct {
	PostID     string    `db:"post_id"`
	CustomerID string    `db:"customer_id"`
	CreatedAt  time.Time `db:"created_at"`
}
