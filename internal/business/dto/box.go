package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type GetNearbyBoxesReq struct {
	Lat        float64 `form:"lat" binding:"required,latitude"`
	Lon        float64 `form:"lon" binding:"required,longitude"`
	Limit      int     `form:"limit" binding:"required,min=1,max=100"`
	Skip       int     `form:"skip" binding:"min=0"`
	CategoryID *int    `form:"category_id" binding:"omitempty,min=1"`
}

type NearbyBoxesResp struct {
	Id            int64           `json:"id"`
	VenueId       int             `json:"venue_id"`
	CategoryID    int             `json:"category_id"`
	Image         string          `json:"image"`
	FullPrice     decimal.Decimal `json:"full_price"`
	DiscountPrice decimal.Decimal `json:"discount_price"`
	CreatedAt     time.Time       `json:"created_at"`
	ExpiresAt     time.Time       `json:"expires_at"`
	Distance      float64         `json:"distance"`
}

type ListResponse struct {
	Items []NearbyBoxesResp
	Total int
}
