package dto

import "time"

type CreateBoxRequest struct {
	VenueID       int       `json:"venue_id" binding:"required,gt=0"`
	CategoryID    *int      `json:"category_id"`
	Image         string    `json:"image" binding:"required"`
	PriceFull     float64   `json:"price_full" binding:"required,gte=0"`
	PriceDiscount float64   `json:"price_discount" binding:"gte=0"`
	ExpiresAt     time.Time `json:"expires_at" binding:"required"`
	Quantity      int       `json:"quantity" binding:"required,gt=0"`
}
