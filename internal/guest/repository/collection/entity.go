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

func executeSQLError(err error) error {
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}

func scanRowsError(err error) error {
	return fmt.Errorf("scan rows: %w", err)
}
