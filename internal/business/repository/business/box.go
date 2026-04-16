package business

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

func (r *Repository) ListNearbyBoxes (ctx context.Context, offset, limit int, lat, lon float64, categoryID *int) (pagination.Result[entity.BoxWithDistance], error) {
	var args []any
	where := `boxes.expires_at > NOW()
			  AND org_units.latitude IS NOT NULL
			  AND org_units.longitude IS NOT NULL
              AND EXISTS (
                  SELECT 1 
                  FROM business.box_items bi 
                  WHERE bi.box_id = boxes.id 
                    AND bi.reserved_by_user_id IS NULL
              )`

	if categoryID != nil {
		args = append(args, *categoryID)
		where += fmt.Sprintf(" AND boxes.category_id=$%d", len(args))
	}

	scanner := func (rows pgx.Rows) (entity.BoxWithDistance, error) {
		var item entity.BoxWithDistance
		err := rows.Scan(
			&item.Box.Id,
			&item.Box.VenueId,
			&item.Box.CategoryID,
			&item.Box.Image,
			&item.Box.FullPrice,
			&item.Box.DiscountPrice,
			&item.Box.CreatedAt,
			&item.Box.ExpiresAt,
			&item.AvailabilityCount,
			&item.Distance,
		)

		if err != nil {
			return entity.BoxWithDistance{}, err
		}
		return item, nil
	}

	dynamicColumns := fmt.Sprintf("boxes.id, boxes.venue_id, boxes.category_id, " +
	"boxes.image, boxes.price_full, boxes.price_discount, " +
	"boxes.created_at, boxes.expires_at, " +
	"(SELECT COUNT(*) FROM business.box_items bi WHERE reserved_by_user_id IS NULL AND bi.box_id=boxes.id) AS availability_count, " +
	"point(%f, %f) <@> point(org_units.longitude, org_units.latitude) AS distance", lon, lat)

	p := pagination.Params{
		Table: "business.org_units JOIN business.boxes on boxes.venue_id=org_units.id",
		Columns: dynamicColumns,
		Where: where,
		OrderBy: "distance ASC, boxes.id ASC",
		Args: args,
		Offset: offset,
		Limit: limit,
	}
	return pagination.List(ctx, r.db.DB(), "business_repository.ListNearbyBoxes", p, scanner)
}