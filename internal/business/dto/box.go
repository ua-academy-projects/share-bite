package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateBoxRequest struct {
	VenueID       int             `json:"venue_id" binding:"required,gt=0"`
	CategoryID    *int            `json:"category_id"`
	Image         string          `json:"image" binding:"required"`
	FullPrice     decimal.Decimal `json:"price_full" binding:"required"`
	DiscountPrice decimal.Decimal `json:"price_discount"`
	ExpiresAt     time.Time       `json:"expires_at" binding:"required"`
	Quantity      int             `json:"quantity" binding:"required,gt=0"`
}

type GetNearbyBoxesReq struct {
	Lat        float64 `form:"lat" binding:"required,latitude"`
	Lon        float64 `form:"lon" binding:"required,longitude"`
	Limit      int     `form:"limit" binding:"min=0,max=100"`
	Skip       int     `form:"skip" binding:"min=0"`
	CategoryID *int    `form:"category_id" binding:"omitempty,min=1"`
}

type NearbyBoxesResp struct {
	Id                 int64           `json:"id" example:"123"`
	VenueID            int             `json:"venue_id" example:"1"`
	CategoryID         *int            `json:"category_id" example:"2"`
	Image              string          `json:"image" example:"https://example.com/box.jpg"`
	FullPrice          decimal.Decimal `json:"full_price" example:"300.00"`
	DiscountPrice      decimal.Decimal `json:"discount_price" example:"150.00"`
	CreatedAt          time.Time       `json:"created_at" example:"2026-04-16T12:00:00Z"`
	ExpiresAt          time.Time       `json:"expires_at" example:"2026-05-01T18:00:00Z"`
	AvailabilityStatus string          `json:"availability_status" example:"running_low"`
	Distance           float64         `json:"distance" example:"1.5"`
}

type ListResponse struct {
	Items []NearbyBoxesResp `json:"items"`
	Total int               `json:"total"`
}
