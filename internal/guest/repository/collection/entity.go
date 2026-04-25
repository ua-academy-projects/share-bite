package collection

import (
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type Collection struct {
	ID         string `db:"id"`
	CustomerID string `db:"customer_id"`

	Name        string  `db:"name"`
	Description *string `db:"description"`
	IsPublic    bool    `db:"is_public"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e Collection) ToEntity() entity.Collection {
	return entity.Collection{
		ID:         e.ID,
		CustomerID: e.CustomerID,

		Name:        e.Name,
		Description: e.Description,
		IsPublic:    e.IsPublic,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type Collections []Collection

func (e Collections) ToEntities() []entity.Collection {
	list := make([]entity.Collection, 0, len(e))
	for _, c := range e {
		list = append(list, c.ToEntity())
	}

	return list
}

type CollectionVenue struct {
	CollectionID string `db:"collection_id"`
	VenueID      int64  `db:"venue_id"`

	SortOrder float64   `db:"sort_order"`
	AddedAt   time.Time `db:"added_at"`
}

func (e CollectionVenue) ToEntity() entity.CollectionVenue {
	return entity.CollectionVenue{
		CollectionID: e.CollectionID,
		VenueID:      e.VenueID,

		SortOrder: e.SortOrder,
		AddedAt:   e.AddedAt,
	}
}

type CollectionVenues []CollectionVenue

func (e CollectionVenues) ToEntities() []entity.CollectionVenue {
	list := make([]entity.CollectionVenue, 0, len(e))
	for _, v := range e {
		list = append(list, v.ToEntity())
	}

	return list
}

type Collaborator struct {
	CollectionID string `db:"collection_id"`
	CustomerID   string `db:"customer_id"`

	UserName        string  `db:"username"`
	AvatarObjectKey *string `db:"avatar_object_key"`

	AddedAt time.Time `db:"added_at"`
}

func (e Collaborator) ToEntity() entity.Collaborator {
	return entity.Collaborator{
		CollectionID: e.CollectionID,
		CustomerID:   e.CustomerID,

		UserName:        e.UserName,
		AvatarObjectKey: e.AvatarObjectKey,

		AddedAt: e.AddedAt,
	}
}

type Collaborators []Collaborator

func (e Collaborators) ToEntities() []entity.Collaborator {
	list := make([]entity.Collaborator, 0, len(e))
	for _, c := range e {
		list = append(list, c.ToEntity())
	}

	return list
}

type Invitation struct {
	ID           string `db:"id"`
	CollectionID string `db:"collection_id"`

	Status entity.InvitationStatus `db:"status"`

	InviterID string `db:"inviter_id"`
	InviteeID string `db:"invitee_id"`

	ExpiresAt  time.Time `db:"expires_at"`
	LastSentAt time.Time `db:"last_sent_at"`
	CreatedAt  time.Time `db:"created_at"`
}

func (e Invitation) ToEntity() entity.Invitation {
	return entity.Invitation{
		ID:           e.ID,
		CollectionID: e.CollectionID,

		Status: e.Status,

		InviterID: e.InviterID,
		InviteeID: e.InviteeID,

		ExpiresAt:  e.ExpiresAt,
		LastSentAt: e.LastSentAt,

		CreatedAt: e.CreatedAt,
	}
}

type EnrichedInvitation struct {
	ID        string                  `db:"id"`
	Status    entity.InvitationStatus `db:"status"`
	CreatedAt time.Time               `db:"created_at"`
	ExpiresAt time.Time               `db:"expires_at"`

	CollectionID   string `db:"collection_id"`
	CollectionName string `db:"collection_name"`

	InviterID              string  `db:"inviter_id"`
	InviterUserName        string  `db:"inviter_username"`
	InviterAvatarObjectKey *string `db:"inviter_avatar_object_key"`

	InviteeID              string  `db:"invitee_id"`
	InviteeUserName        string  `db:"invitee_username"`
	InviteeAvatarObjectKey *string `db:"invitee_avatar_object_key"`
}

func (e EnrichedInvitation) ToEntity() entity.EnrichedInvitation {
	return entity.EnrichedInvitation{
		ID:        e.ID,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		ExpiresAt: e.ExpiresAt,

		CollectionID:   e.CollectionID,
		CollectionName: e.CollectionName,

		InviterID:              e.InviterID,
		InviterUserName:        e.InviterUserName,
		InviterAvatarObjectKey: e.InviterAvatarObjectKey,

		InviteeID:              e.InviteeID,
		InviteeUserName:        e.InviteeUserName,
		InviteeAvatarObjectKey: e.InviteeAvatarObjectKey,
	}
}

type EnrichedInvitations []EnrichedInvitation

func (e EnrichedInvitations) ToEntities() []entity.EnrichedInvitation {
	list := make([]entity.EnrichedInvitation, 0, len(e))
	for _, ei := range e {
		list = append(list, ei.ToEntity())
	}

	return list
}

func executeSQLError(err error) error {
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}

func scanRowsError(err error) error {
	return fmt.Errorf("scan rows: %w", err)
}
