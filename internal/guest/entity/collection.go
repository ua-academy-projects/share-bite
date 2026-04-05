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
	VenueID      string
	SortOrder    float64
	AddedAt      time.Time
}

// Venue describes the venue FROM THE PERSPECTIVE OF the collections module.
type Venue struct {
	ID          string
	Name        string
	Description *string

	AvatarURL *string
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
	PageSize   int
	PageToken  string
}

type ListCustomerCollectionsOutput struct {
	Collections   []Collection
	NextPageToken string
}

//

type ReorderVenueInput struct {
	CollectionID string
	VenueID      string

	CustomerID string

	PrevVenueID *string
	NextVenueID *string
}
