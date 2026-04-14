package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Box struct {
	Id            int64
	VenueId       int
	CategoryID    int
	Image         string
	FullPrice     decimal.Decimal
	DiscountPrice decimal.Decimal
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

type BoxWithDistance struct {
    Box      Box
    Distance float64
}