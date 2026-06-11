package business

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

func (r *Repository) GetDailySummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID) (entity.DailySummary, error) {
	const op = "repository.business.GetDailySummary"

	q := database.Query{
		Name: "GetDailySummary",
		Sql: `WITH brand AS (
				SELECT id 
				FROM business.org_units
				WHERE org_account_id = $1 AND profile_type = 'BRAND'
			),
			venues AS (
				SELECT id 
				FROM business.org_units
				WHERE parent_id = (SELECT id FROM brand) AND profile_type = 'VENUE'
			)
			SELECT 
				(SELECT COUNT(*) FROM venues) AS total_venues_count,
				
				(SELECT COUNT(*) 
				 FROM business.boxes 
				 WHERE venue_id IN (SELECT id FROM venues)
				   AND created_at >= $2 AND created_at <= $3) AS created_boxes_count,
				   
				(SELECT COUNT(*) 
				 FROM business.posts 
				 WHERE org_id IN (SELECT id FROM brand UNION SELECT id FROM venues)
				   AND created_at >= $2 AND created_at <= $3) AS created_posts_count;
		`,
	}

	var sum entity.DailySummary

	err := r.db.DB().QueryRowContext(ctx, q, orgID, startDate, endDate).Scan(
		&sum.TotalVenuesCount,
		&sum.CreatedBoxesCount,
		&sum.CreatedPostsCount,
	)

	if err != nil {
		return entity.DailySummary{}, err
	}
	return sum, nil
}

func (r *Repository) GetReservationSummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.ReservationSummary, error) {
	const op = "repository.business.GetReservationSummary"

	q := database.Query{
		Name: "GetReservationSummary",
		Sql: `WITH target_venues AS (
				SELECT id 
				FROM business.org_units
					WHERE ($4::int IS NOT NULL AND id = $4)
					OR ($4::int IS NULL AND parent_id = (
						SELECT id FROM business.org_units
						WHERE org_account_id = $3 AND profile_type = 'BRAND'
					))
				)
				SELECT 
					COUNT(bi.box_code) FILTER (WHERE bi.status = 'SOLD') AS sold_items,
					COUNT(bi.box_code) FILTER (WHERE bi.reserved_by_user_id IS NOT NULL AND bi.status != 'SOLD') AS reserved_items,
					COUNT(bi.box_code) FILTER (WHERE bi.reserved_by_user_id IS NULL) AS available_items,
					COALESCE(SUM(b.price_discount) FILTER (WHERE bi.reserved_by_user_id IS NOT NULL), 0) AS potential_revenue
				FROM business.boxes b
				JOIN business.box_items bi ON b.id = bi.box_id
				WHERE b.venue_id IN (SELECT id FROM target_venues)
				AND b.created_at >= $1 
				AND b.created_at <= $2;`,
	}

	var res entity.ReservationSummary
	err := r.db.DB().QueryRowContext(ctx, q, startDate, endDate, orgID, venueID).Scan(
		&res.TotalSoldItems,
		&res.TotalReservedItems,
		&res.TotalAvailableItems,
		&res.PotentialRevenue,
	)

	if err != nil {
		return entity.ReservationSummary{}, err
	}
	return res, nil
}

func (r *Repository) GetVenueActivitySummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID int) (entity.VenueActivitySummary, error) {
	const op = "repository.business.GetVenueActivitySummary"
	q := database.Query{
		Name: "GetVenueActivity",
		Sql: `WITH venue_info AS (
   				SELECT name
    			FROM business.org_units
    			WHERE id = $4
			)
			SELECT 
    			(SELECT COUNT(*) FROM business.boxes 
     			WHERE venue_id = $4 
       			AND created_at >= $1 AND created_at <= $2) AS total_boxes_created,
       
    			(SELECT COUNT(*) FROM business.posts 
     			WHERE org_id = $4 
       			AND created_at >= $1 AND created_at <= $2) AS total_posts_created,
       
    			(SELECT name FROM venue_info) AS venue_name;`,
	}

	var res entity.VenueActivitySummary
	err := r.db.DB().QueryRowContext(ctx, q, startDate, endDate, orgID, venueID).Scan(
		&res.TotalBoxesCreated,
		&res.TotalPostsCreated,
		&res.VenueName,
	)
	if err != nil {
		return entity.VenueActivitySummary{}, err
	}
	return res, nil
}

func (r *Repository) GetFoodBoxPerformance(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.BoxPerformanceRaw, error) {
	const op = "repository.business.GetFoodBoxPerformance"
	q := database.Query{
		Name: "GetFoodBoxPerformance",
		Sql: `WITH target_venues AS (
    			SELECT id
    			FROM business.org_units
    			WHERE ($4::int IS NOT NULL AND id = $4)
       			OR ($4::int IS NULL AND parent_id = (
           			SELECT id FROM business.org_units 
           			WHERE org_account_id = $3 AND profile_type = 'BRAND'
       			))
			),
			target_boxes AS (
				SELECT id, price_full, price_discount, expires_at
				FROM business.boxes
				WHERE venue_id IN (SELECT id FROM target_venues)
				AND created_at >= $1 
				AND created_at <= $2
			)
			SELECT 
				(SELECT COUNT(*) FROM target_boxes) AS total_boxes_created,
				
				(SELECT COUNT(*) FROM target_boxes WHERE expires_at < NOW()) AS total_boxes_expired,
				
				(SELECT COALESCE(AVG(price_full - price_discount), 0) FROM target_boxes) AS average_discount,
				
				(SELECT COUNT(*) FROM business.box_items 
				WHERE box_id IN (SELECT id FROM target_boxes)) AS total_box_items,
				
				(SELECT COUNT(*) FROM business.box_items 
				WHERE box_id IN (SELECT id FROM target_boxes) 
				AND reserved_by_user_id IS NOT NULL) AS total_reserved_items;`,
	}

	var res entity.BoxPerformanceRaw
	err := r.db.DB().QueryRowContext(ctx, q, startDate, endDate, orgID, venueID).Scan(
		&res.TotalBoxesCreated,
		&res.TotalBoxesExpired,
		&res.AverageDiscount,
		&res.TotalBoxItems,
		&res.TotalReservedItems,
	)
	if err != nil {
		return entity.BoxPerformanceRaw{}, err
	}

	return res, nil
}

func (r *Repository) GetEngagementSummary(ctx context.Context, startDate, endDate time.Time, orgID uuid.UUID, venueID *int) (entity.EngagementSummaryRaw, error) {
	const op = "repository.business.GetEngagementSummary"
	q := database.Query{
		Name: "GetEngagementSummary",
		Sql: `WITH target_venues AS (
				SELECT id
				FROM business.org_units
				WHERE ($4::int IS NOT NULL AND id = $4)
				OR ($4::int IS NULL AND parent_id = (
					SELECT id FROM business.org_units 
					WHERE org_account_id = $3 AND profile_type = 'BRAND'
				))
			),
			target_posts AS (
				SELECT id
				FROM business.posts
				WHERE org_id IN (SELECT id FROM target_venues)
				AND created_at >= $1 
				AND created_at <= $2
			)
			SELECT 
				(SELECT COUNT(*) FROM target_posts) AS total_posts_created,
				
				(SELECT COUNT(*) FROM business.comments 
				WHERE post_id IN (SELECT id FROM target_posts)) AS total_comments,
				
				(SELECT COUNT(*) FROM business.likes 
				WHERE post_id IN (SELECT id FROM target_posts)) AS total_likes;`,
	}

	var res entity.EngagementSummaryRaw
	err := r.db.DB().QueryRowContext(ctx, q, startDate, endDate, orgID, venueID).Scan(
		&res.TotalPostsCreated,
		&res.TotalComments,
		&res.TotalLikes,
	)
	if err != nil {
		return entity.EngagementSummaryRaw{}, err
	}
	return res, nil
}
