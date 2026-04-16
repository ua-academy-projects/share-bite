package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

const RunningLowThreshold = 7

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
	Box               Box
	AvailabilityCount int
	Distance          float64
}

func (b BoxWithDistance) AvailabilityStatus() string{
	if b.AvailabilityCount <= RunningLowThreshold{
		return "running_low"
	}else {
		return "available"
	}
}