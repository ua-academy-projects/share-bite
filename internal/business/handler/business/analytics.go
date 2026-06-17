package business

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/middleware"

	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
)

func parseDates(startDate, endDate string) (*time.Time, *time.Time, error) {
	parsedStartDate, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return nil, nil, err
	}

	parsedEndDate, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return nil, nil, err
	}
	return &parsedStartDate, &parsedEndDate, nil
}

func isStartDateValid(startDate, endDate time.Time) bool {
	if startDate.After(time.Now()) {
		return false
	}
	if startDate.After(endDate) {
		return false
	}
	if endDate.Sub(startDate) > 90*time.Hour*24 {
		return false
	}
	return true
}

func (h *handler) GetDailySummary(c *gin.Context) {
	var req dto.DailySummaryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(apperror.BadRequest("Invalid request"))
		return
	}

	startDate, endDate, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing dates"))
		return
	}

	if !isStartDateValid(*startDate, *endDate) {
		c.Error(apperror.BadRequest("Invalid start date"))
		return
	}

	userIDStr, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized access attempt"))
		return
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing id"))
		return
	}

	res, err := h.service.GetDailySummary(c, *startDate, *endDate, parsedID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.DailySummaryResponse{
		CreatedBoxesCount: res.CreatedBoxesCount,
		CreatedPostsCount: res.CreatedPostsCount,
		TotalVenuesCount:  res.TotalVenuesCount,
	})
}

func (h *handler) GetReservationSummary(c *gin.Context) {
	var req dto.ReservationSummaryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(apperror.BadRequest("Invalid request"))
		return
	}

	startDate, endDate, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing dates"))
		return
	}

	if !isStartDateValid(*startDate, *endDate) {
		c.Error(apperror.BadRequest("Invalid start date"))
		return
	}

	userIDStr, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized access attempt"))
		return
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing id"))
		return
	}

	res, err := h.service.GetReservationSummary(c, *startDate, *endDate, parsedID, req.VenueID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ReservationSummaryResponse{
		TotalSoldItems:      res.TotalSoldItems,
		TotalReservedItems:  res.TotalReservedItems,
		TotalAvailableItems: res.TotalAvailableItems,
		PotentialRevenue:    res.PotentialRevenue,
	})
}

func (h *handler) GetVenueActivitySummary(c *gin.Context) {
	var req dto.VenueActivitySummaryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(apperror.BadRequest("Invalid request"))
		return
	}

	startDate, endDate, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing dates"))
		return
	}

	if !isStartDateValid(*startDate, *endDate) {
		c.Error(apperror.BadRequest("Invalid start date"))
		return
	}

	userIDStr, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized access attempt"))
		return
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing id"))
		return
	}

	venueIDStr := c.Param("venue_id")
	venueID, err := strconv.Atoi(venueIDStr)
	if err != nil {
		c.Error(apperror.BadRequest("Invalid venue_id"))
		return
	}

	res, err := h.service.GetVenueActivitySummary(c, *startDate, *endDate, parsedID, venueID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.VenueActivitySummaryResponse{
		TotalBoxesCreated: res.TotalBoxesCreated,
		TotalPostsCreated: res.TotalPostsCreated,
		VenueName:         res.VenueName,
	})
}

func (h *handler) GetFoodBoxPerformance(c *gin.Context) {
	var req dto.FoodBoxPerformanceRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(apperror.BadRequest("Invalid request"))
		return
	}

	startDate, endDate, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing dates"))
		return
	}

	if !isStartDateValid(*startDate, *endDate) {
		c.Error(apperror.BadRequest("Invalid start date"))
		return
	}

	userIDStr, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized access attempt"))
		return
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing id"))
		return
	}

	res, err := h.service.GetFoodBoxPerformance(c, *startDate, *endDate, parsedID, req.VenueID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.FoodBoxPerformanceResponse{
		TotalBoxesCreated: res.TotalBoxesCreated,
		TotalBoxesExpired: res.TotalBoxesExpired,
		AverageDiscount:   res.AverageDiscount,
		SellThroughRate:   res.SellThroughRate,
		WasteRate:         res.WasteRate,
		Score:             res.Score,
	})
}

func (h *handler) GetEngagementSummary(c *gin.Context) {
	var req dto.EngagementSummaryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request"))
		return
	}

	startDate, endDate, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing dates"))
		return
	}

	if !isStartDateValid(*startDate, *endDate) {
		c.Error(apperror.BadRequest("Invalid start date"))
		return
	}

	userIDStr, ok := middleware.GetUserID(c)
	if !ok {
		c.Error(apperror.Unauthorized("unauthorized access attempt"))
		return
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Error(apperror.BadRequest("Error parsing id"))
		return
	}

	res, err := h.service.GetEngagementSummary(c, *startDate, *endDate, parsedID, req.VenueID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.EngagementSummaryResponse{
		TotalPostsCreated:  res.TotalPostsCreated,
		TotalComments:      res.TotalComments,
		TotalLikes:         res.TotalLikes,
		AverageCommentsNum: res.AverageCommentsNum,
		AverageLikesNum:    res.AverageLikesNum,
	})
}
