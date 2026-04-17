package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Box struct {
	ID            int64
	VenueID       int
	CategoryID    *int

	Image         string
	PriceFull     decimal.Decimal
	PriceDiscount decimal.Decimal

	CreatedAt time.Time
	ExpiresAt time.Time
}

type BoxItem struct {
	BoxID            int64
	BoxCode          string
	ReservedByUserID *string
}

type BoxWithDistance struct {
	Box      Box
	Distance float64
}