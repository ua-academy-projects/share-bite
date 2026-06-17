package dto

import "github.com/shopspring/decimal"

type DailySummaryResponse struct {
	CreatedBoxesCount int `json:"created_boxes_count"`
	CreatedPostsCount int `json:"created_posts_count"`
	TotalVenuesCount  int `json:"total_venues_count"`
}

type DailySummaryRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type ReservationSummaryResponse struct {
	TotalSoldItems      int             `json:"total_sold_items"`
	TotalReservedItems  int             `json:"total_reserved_items"`
	TotalAvailableItems int             `json:"total_available_items"`
	PotentialRevenue    decimal.Decimal `json:"potential_revenue"` // sum of discounted prices of reserved box items
}

type ReservationSummaryRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	VenueID   *int   `form:"venue_id"`
}

type VenueActivitySummaryResponse struct {
	TotalBoxesCreated int    `json:"total_boxes_created"`
	TotalPostsCreated int    `json:"total_posts_created"`
	VenueName         string `json:"venue_name"`
}

type VenueActivitySummaryRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

type FoodBoxPerformanceResponse struct {
	TotalBoxesCreated int     `json:"total_boxes_created"`
	TotalBoxesExpired int     `json:"total_boxes_expired"`
	AverageDiscount   float64 `json:"average_discount"`
	SellThroughRate   float64 `json:"sell_through_rate"` // reserved/total
	WasteRate         float64 `json:"waste_rate"`        // expired/total
	Score             float64 `json:"score"`             // (selfThroughRate * 100 * 0.7) + ((1 - wasteRate) * 100 * 0.3)
}

type FoodBoxPerformanceRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	VenueID   *int   `form:"venue_id"`
}

type EngagementSummaryResponse struct {
	TotalPostsCreated  int     `json:"total_posts_created"`
	TotalComments      int     `json:"total_comments"`
	TotalLikes         int     `json:"total_likes"`
	AverageCommentsNum float64 `json:"average_comments_num"`
	AverageLikesNum    float64 `json:"average_likes_num"`
}

type EngagementSummaryRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	VenueID   *int   `form:"venue_id"`
}
