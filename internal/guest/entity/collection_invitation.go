package entity

import (
	"time"
)

type InvitationStatus string

const (
	PendingInvitationStatus  InvitationStatus = "pending"
	AcceptedInvitationStatus InvitationStatus = "accepted"
	DeclinedInvitationStatus InvitationStatus = "declined"
)

type Invitation struct {
	ID           string
	CollectionID string

	Status InvitationStatus

	InviterID string
	InviteeID string

	CreatedAt  time.Time
	LastSentAt time.Time
	ExpiresAt  time.Time
}

func (i Invitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

func (i Invitation) CanResend(cooldown time.Duration) bool {
	return time.Since(i.LastSentAt) >= cooldown
}

func (i Invitation) CanBeAccepted() bool {
	return i.Status == PendingInvitationStatus && !i.IsExpired()
}

func (i Invitation) CanBeDeclined() bool {
	return i.Status == PendingInvitationStatus
}

// EnrichedInvitation represents a read-optimized model of an Invitation.
// It aggregates data from the invitations, collections, and customers tables.
type EnrichedInvitation struct {
	ID        string
	Status    InvitationStatus
	CreatedAt time.Time
	ExpiresAt time.Time

	CollectionID   string
	CollectionName string

	InviterID              string
	InviterUserName        string
	InviterAvatarObjectKey *string

	InviteeID              string
	InviteeUserName        string
	InviteeAvatarObjectKey *string
}

type InviteCollaboratorInput struct {
	CollectionID string

	InviterID string
	InviteeID string

	Expiry time.Time
}

type ListInvitationsInput struct {
	CollectionID *string

	InviterID *string
	InviteeID *string

	CallerID string

	Status *InvitationStatus

	CursorID string
	Limit    int
}

type ListInvitationsOutput struct {
	Invitations  []EnrichedInvitation
	NextCursorID string
}
