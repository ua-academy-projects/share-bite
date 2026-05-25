package business

import (
	"context"
	stdsql "database/sql"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

func (r *Repository) ReplaceLocationHours(ctx context.Context, venueID int, days []dto.VenueHoursDayInput) error {
	delQ := database.Query{
		Name: "business_repository.DeleteLocationHours",
		Sql:  `DELETE FROM business.location_hours WHERE venue_id = $1`,
	}
	if _, err := r.db.DB().ExecContext(ctx, delQ, venueID); err != nil {
		return executeSQLError(err)
	}

	insQ := database.Query{
		Name: "business_repository.InsertLocationHours",
		Sql: `
			INSERT INTO business.location_hours (venue_id, weekday, open_time, close_time)
			VALUES ($1, $2, $3::time, $4::time)
		`,
	}

	for _, d := range days {
		if _, err := r.db.DB().ExecContext(ctx, insQ, venueID, d.Weekday, d.OpenTime, d.CloseTime); err != nil {
			return executeSQLError(err)
		}
	}

	return nil
}

func (r *Repository) GetLocationHours(ctx context.Context, venueID int) ([]dto.VenueHoursDayInput, error) {
	q := database.Query{
		Name: "business_repository.GetLocationHours",
		Sql: `
			SELECT
				weekday,
				TO_CHAR(open_time, 'HH24:MI') AS open_time,
				TO_CHAR(close_time, 'HH24:MI') AS close_time
			FROM business.location_hours
			WHERE venue_id = $1
			ORDER BY weekday
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, venueID)
	if err != nil {
		return nil, fmt.Errorf("query location hours: %w", err)
	}
	defer rows.Close()

	out := make([]dto.VenueHoursDayInput, 0, 7)
	for rows.Next() {
		var (
			weekday   int
			openTime  stdsql.NullString
			closeTime stdsql.NullString
		)
		if err := rows.Scan(&weekday, &openTime, &closeTime); err != nil {
			return nil, fmt.Errorf("scan location hours: %w", err)
		}

		var openPtr, closePtr *string
		if openTime.Valid {
			v := openTime.String
			openPtr = &v
		}
		if closeTime.Valid {
			v := closeTime.String
			closePtr = &v
		}

		out = append(out, dto.VenueHoursDayInput{
			Weekday:   weekday,
			OpenTime:  openPtr,
			CloseTime: closePtr,
		})
	}

	return out, nil
}
