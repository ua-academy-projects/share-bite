package entity

import "time"

type PostCollaboratorStatus string

const (
	PostCollaboratorStatusPending  PostCollaboratorStatus = "pending"
	PostCollaboratorStatusAccepted PostCollaboratorStatus = "accepted"
	PostCollaboratorStatusDeclined PostCollaboratorStatus = "declined"
)

type PostCollaborator struct {
	ID          int64
	PostID      int64
	CustomerID  string
	InvitedBy   string
	Status      PostCollaboratorStatus
	ExpiresAt   time.Time
	RespondedAt *time.Time
	CreatedAt   time.Time
}
