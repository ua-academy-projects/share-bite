package entity

import "github.com/shopspring/decimal"

type DailySummary struct {
	CreatedBoxesCount int
	CreatedPostsCount int
	TotalVenuesCount  int
}

type ReservationSummary struct {
	TotalSoldItems      int
	TotalReservedItems  int
	TotalAvailableItems int
	PotentialRevenue    decimal.Decimal
}

type VenueActivitySummary struct {
	TotalBoxesCreated int
	TotalPostsCreated int
	VenueName         string
}

type BoxPerformanceRaw struct {
	TotalBoxesCreated  int
	TotalBoxesExpired  int
	AverageDiscount    float64
	TotalBoxItems      int
	TotalReservedItems int
}

type BoxPerformance struct {
	TotalBoxesCreated int
	TotalBoxesExpired int
	AverageDiscount   float64
	SellThroughRate   float64
	WasteRate         float64
	Score             float64
}

type EngagementSummaryRaw struct {
	TotalPostsCreated int
	TotalComments     int
	TotalLikes        int
}

type EngagementSummary struct {
	TotalPostsCreated  int
	TotalComments      int
	TotalLikes         int
	AverageCommentsNum float64
	AverageLikesNum    float64
}
