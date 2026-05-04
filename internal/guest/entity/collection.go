package entity

import "time"

type Collection struct {
	ID         string
	CustomerID string

	Name        string
	Description *string
	IsPublic    bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CollectionVenue struct {
	CollectionID string
	VenueID      int64
	SortOrder    float64
	AddedAt      time.Time
}

// Venue describes the venue FROM THE PERSPECTIVE OF the collections module.
type Venue struct {
	ID          int64
	Name        string
	Description *string

	AvatarURL *string
	BannerURL *string
}

type EnrichedVenueItem struct {
	VenueItem Venue

	SortOrder float64
	AddedAt   time.Time
}

//

type CreateCollectionInput struct {
	CustomerID string

	Name        string
	Description *string
	IsPublic    bool
}

type UpdateCollectionInput struct {
	CollectionID string
	CustomerID   string

	Name        *string
	Description *string
	IsPublic    *bool
}

type ListCustomerCollectionsInput struct {
	CustomerID string

	CursorTime time.Time
	CursorID   string
	Limit      int
}

type ListCustomerCollectionsOutput struct {
	Collections    []Collection
	NextCursorTime *time.Time
	NextCursorID   *string
}

//

type ReorderVenueInput struct {
	CollectionID string
	VenueID      int64

	CustomerID string

	PrevVenueID *int64
	NextVenueID *int64
}
