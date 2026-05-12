package post

import "time"

type CollaboratorStatus string

type PostCollaborator struct {
	ID          int64              `db:"id"`
	PostID      int64              `db:"post_id"`
	CustomerID  string             `db:"customer_id"`
	InvitedBy   string             `db:"invited_by"`
	Status      CollaboratorStatus `db:"status"`
	ExpiresAt   time.Time          `db:"expires_at"`
	RespondedAt *time.Time         `db:"responded_at"`
	CreatedAt   time.Time          `db:"created_at"`
}
