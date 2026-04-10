package entity

import "time"

type Box struct {
	ID            int64
	VenueID       int
	CategoryID    *int
	Image         string
	PriceFull     float64
	PriceDiscount float64
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

type BoxItem struct {
	BoxID            int64
	BoxCode          string
	ReservedByUserID *string
}
