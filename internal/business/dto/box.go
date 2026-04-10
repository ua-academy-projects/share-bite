package dto

import "time"

type CreateBoxRequest struct {
	VenueID       int       `json:"venue_id"`
	CategoryID    *int      `json:"category_id"`
	Image         string    `json:"image"`
	PriceFull     float64   `json:"price_full"`
	PriceDiscount float64   `json:"price_discount"`
	ExpiresAt     time.Time `json:"expires_at"`
	Quantity      int       `json:"quantity"`
}
