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
	ID            int64           `json:"id"`
	VenueID       int             `json:"venue_id"`
	CategoryID    *int             `json:"category_id"`
	Image         string          `json:"image"`
	FullPrice     decimal.Decimal `json:"full_price"`
	DiscountPrice decimal.Decimal `json:"discount_price"`
	CreatedAt     time.Time       `json:"created_at"`
	ExpiresAt     time.Time       `json:"expires_at"`
	Distance      float64         `json:"distance"`
}

type ListResponse struct {
	Items []NearbyBoxesResp `json:"items"`
	Total int               `json:"total"`
}