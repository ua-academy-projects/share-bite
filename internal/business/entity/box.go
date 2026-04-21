package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

const RunningLowThreshold = 7

type Box struct {
	ID         int64
	VenueID    int
	CategoryID *int

	Image         string
	FullPrice     decimal.Decimal
	DiscountPrice decimal.Decimal

	CreatedAt time.Time
	ExpiresAt time.Time
}

type BoxItem struct {
	BoxID            int64
	BoxCode          string
	ReservedByUserID *string
}

type BoxWithDistance struct {
	Box               Box
	AvailabilityCount int
	Distance          float64
}

func (b BoxWithDistance) AvailabilityStatus() string {
	if b.AvailabilityCount == 0 {
		return "sold_out"
	} else if b.AvailabilityCount <= RunningLowThreshold {
		return "running_low"
	} else {
		return "available"
	}
}
