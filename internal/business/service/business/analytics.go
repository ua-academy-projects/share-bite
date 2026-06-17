package business

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

var (
	sellCoef  = 0.7
	wasteCoef = 0.3
)

func (s *service) GetDailySummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID) (entity.DailySummary, error) {
	result, err := s.businessRepo.GetDailySummary(ctx, startDate, endDate, orgID)
	if err != nil {
		return entity.DailySummary{}, fmt.Errorf("Database error: %w", err)
	}
	return result, nil
}

func (s *service) GetReservationSummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.ReservationSummary, error) {
	if venueID != nil {
		if err := s.businessRepo.CheckOwnership(ctx, orgID.String(), *venueID); err != nil {
			return entity.ReservationSummary{}, err
		}
	}

	result, err := s.businessRepo.GetReservationSummary(ctx, startDate, endDate, orgID, venueID)
	if err != nil {
		return entity.ReservationSummary{}, fmt.Errorf("Database error: %w", err)
	}
	return result, nil
}

func (s *service) GetVenueActivitySummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID int) (entity.VenueActivitySummary, error) {
	if err := s.businessRepo.CheckOwnership(ctx, orgID.String(), venueID); err != nil {
		return entity.VenueActivitySummary{}, err
	}

	result, err := s.businessRepo.GetVenueActivitySummary(ctx, startDate, endDate, orgID, venueID)
	if err != nil {
		return entity.VenueActivitySummary{}, err
	}
	return result, nil
}

func (s *service) GetFoodBoxPerformance(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.BoxPerformance, error) {
	if venueID != nil {
		if err := s.businessRepo.CheckOwnership(ctx, orgID.String(), *venueID); err != nil {
			return entity.BoxPerformance{}, err
		}
	}

	raw, err := s.businessRepo.GetFoodBoxPerformance(ctx, startDate, endDate, orgID, venueID)
	if err != nil {
		return entity.BoxPerformance{}, err
	}

	sellThroughRate := 0.0
	if raw.TotalBoxItems > 0 {
		sellThroughRate = float64(raw.TotalReservedItems) / float64(raw.TotalBoxItems)
	}
	wasteRate := 0.0
	if raw.TotalBoxesCreated > 0 {
		wasteRate = float64(raw.TotalBoxesExpired) / float64(raw.TotalBoxesCreated)
	}

	score := (sellThroughRate * 100.0 * sellCoef) + ((1 - wasteRate) * 100.0 * wasteCoef)

	return entity.BoxPerformance{
		TotalBoxesCreated: raw.TotalBoxesCreated,
		TotalBoxesExpired: raw.TotalBoxesExpired,
		AverageDiscount:   raw.AverageDiscount,
		SellThroughRate:   sellThroughRate,
		WasteRate:         wasteRate,
		Score:             score,
	}, nil
}

func (s *service) GetEngagementSummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.EngagementSummary, error) {
	if venueID != nil {
		if err := s.businessRepo.CheckOwnership(ctx, orgID.String(), *venueID); err != nil {
			return entity.EngagementSummary{}, err
		}
	}

	raw, err := s.businessRepo.GetEngagementSummary(ctx, startDate, endDate, orgID, venueID)
	if err != nil {
		return entity.EngagementSummary{}, err
	}

	var averageCommentsNum, averageLikesNum float64
	if raw.TotalPostsCreated > 0 {
		averageCommentsNum = float64(raw.TotalComments) / float64(raw.TotalPostsCreated)
		averageLikesNum = float64(raw.TotalLikes) / float64(raw.TotalPostsCreated)
	}

	return entity.EngagementSummary{
		TotalPostsCreated:  raw.TotalPostsCreated,
		TotalComments:      raw.TotalComments,
		TotalLikes:         raw.TotalLikes,
		AverageCommentsNum: averageCommentsNum,
		AverageLikesNum:    averageLikesNum,
	}, nil
}
