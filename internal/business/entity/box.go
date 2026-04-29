package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

const RunningLowThreshold = 7

type Status string

const (
	StatusSoldOut    Status = "sold_out"
	StatusRunningLow Status = "running_low"
	StatusAvailable  Status = "available"
)

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

func (b BoxWithDistance) AvailabilityStatus() Status {
	if b.AvailabilityCount == 0 {
		return StatusSoldOut
	} else if b.AvailabilityCount <= RunningLowThreshold {
		return StatusRunningLow
	} else {
		return StatusAvailable
	}
}
