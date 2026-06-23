package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
)

func (r *Repository) CreateBox(ctx context.Context, box *entity.Box) (int64, time.Time, error) {
	const op = "repository.box.CreateBox"

	q := database.Query{
		Name: "create_box",
		Sql: `
		INSERT INTO business.boxes 
        (venue_id, category_id, image, price_full, price_discount, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
	`,
	}

	var id int64
	var createdAt time.Time

	if tx, ok := ctx.Value(pg.TxKey).(pgx.Tx); ok {
		err := tx.QueryRow(ctx, q.Sql,
			box.VenueID,
			box.CategoryID,
			box.Image,
			box.FullPrice,
			box.DiscountPrice,
			box.ExpiresAt,
		).Scan(&id, &createdAt)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
				return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrNotFound)
			}
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
		}

		return id, createdAt, nil
	}

	err := r.db.DB().QueryRowContext(ctx, q,
		box.VenueID,
		box.CategoryID,
		box.Image,
		box.FullPrice,
		box.DiscountPrice,
		box.ExpiresAt,
	).Scan(&id, &createdAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, createdAt, nil
}

func (r *Repository) CreateBoxItem(ctx context.Context, boxID int64, code string) error {
	const op = "repository.box.CreateBoxItem"

	q := database.Query{
		Name: "create_box_item",
		Sql: `
		INSERT INTO business.box_items (box_id, box_code)
        VALUES ($1, $2)
	`,
	}

	if tx, ok := ctx.Value(pg.TxKey).(pgx.Tx); ok {
		_, err := tx.Exec(ctx, q.Sql, boxID, code)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	_, err := r.db.DB().ExecContext(ctx, q, boxID, code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int, orgID *int) (pagination.Result[entity.BoxWithDistance], error) {
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
	if orgID != nil {
		args = append(args, *orgID)
		where += fmt.Sprintf(" AND (org_units.id=$%d OR org_units.parent_id=$%d)", len(args), len(args))
	}

	scanner := func(rows pgx.Rows) (entity.BoxWithDistance, error) {
		var item entity.BoxWithDistance
		var venueStatus string

		err := rows.Scan(
			&item.Box.ID,
			&item.Box.VenueID,
			&item.Box.CategoryID,
			&item.Box.Image,
			&item.Box.FullPrice,
			&item.Box.DiscountPrice,
			&item.Box.CreatedAt,
			&item.Box.ExpiresAt,
			&item.AvailabilityCount,
			&item.Distance,
			&venueStatus,
		)

		if err != nil {
			return entity.BoxWithDistance{}, err
		}

		item.VenueStatus = entity.OrgStatus(venueStatus)

		return item, nil
	}

	dynamicColumns := fmt.Sprintf("boxes.id, boxes.venue_id, boxes.category_id, "+
		"boxes.image, boxes.price_full, boxes.price_discount, "+
		"boxes.created_at, boxes.expires_at, "+
		"(SELECT COUNT(*) FROM business.box_items bi WHERE reserved_by_user_id IS NULL AND bi.box_id=boxes.id) AS availability_count, "+
		"point(%f, %f) <@> point(org_units.longitude, org_units.latitude) AS distance, "+
		"org_units.status",
		lon, lat)

	p := pagination.Params{
		Table:   "business.org_units JOIN business.boxes on boxes.venue_id=org_units.id",
		Columns: dynamicColumns,
		Where:   where,
		OrderBy: "distance ASC, boxes.id ASC",
		Args:    args,
		Offset:  offset,
		Limit:   limit,
	}

	return pagination.List(ctx, r.db.DB(), "business_repository.ListNearbyBoxes", p, scanner)
}

func (r *Repository) GetBox(ctx context.Context, boxID int64) (*entity.Box, error) {
	const op = "repository.box.GetBox"

	q := database.Query{
		Name: "get_box",
		Sql: `
		SELECT id, venue_id, category_id, image, price_full, price_discount, created_at, expires_at
		FROM business.boxes
		WHERE id = $1
	`,
	}

	var box entity.Box
	err := r.db.DB().QueryRowContext(ctx, q, boxID).Scan(
		&box.ID,
		&box.VenueID,
		&box.CategoryID,
		&box.Image,
		&box.FullPrice,
		&box.DiscountPrice,
		&box.CreatedAt,
		&box.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &box, nil
}

func (r *Repository) ReserveBoxItem(ctx context.Context, boxID int64, userID string) (string, error) {
	const op = "repository.box.ReserveBoxItem"

	q := database.Query{
		Name: "reserve_box_item",
		Sql: `
		UPDATE business.box_items
		SET reserved_by_user_id = $1
		WHERE box_code = (
			SELECT box_code
			FROM business.box_items
			WHERE box_id = $2 AND reserved_by_user_id IS NULL
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING box_code
	`,
	}

	var boxCode string
	err := r.db.DB().QueryRowContext(ctx, q, userID, boxID).Scan(&boxCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, ErrNoAvailableItems)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return boxCode, nil
}

func (r *Repository) ListBoxesByVenueID(ctx context.Context, venueID int, offset, limit int) (pagination.Result[entity.Box], error) {
	const op = "repository.box.ListBoxesByVenueID"

	scanner := func(rows pgx.Rows) (entity.Box, error) {
		var box entity.Box

		err := rows.Scan(
			&box.ID,
			&box.VenueID,
			&box.CategoryID,
			&box.Image,
			&box.FullPrice,
			&box.DiscountPrice,
			&box.CreatedAt,
			&box.ExpiresAt,
		)

		if err != nil {
			return entity.Box{}, err
		}

		return box, nil
	}

	p := pagination.Params{
		Table:   "business.boxes",
		Columns: "id, venue_id, category_id, image, price_full, price_discount, created_at, expires_at",
		Where:   "venue_id = $1",
		Args:    []any{venueID},
		OrderBy: "created_at DESC",
		Offset:  offset,
		Limit:   limit,
	}

	result, err := pagination.List(ctx, r.db.DB(), op, p, scanner)
	if err != nil {
		return pagination.Result[entity.Box]{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r *Repository) UpdateBox(ctx context.Context, boxID int64, input entity.BoxUpdateInput) (*entity.Box, error) {
	const op = "repository.box.UpdateBox"

	var setClauses []string
	var args []any
	argNum := 1

	if input.CategoryID != nil {
		setClauses = append(setClauses, fmt.Sprintf("category_id = $%d", argNum))
		args = append(args, *input.CategoryID)
		argNum++
	}

	if input.FullPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("price_full = $%d", argNum))
		args = append(args, *input.FullPrice)
		argNum++
	}

	if input.DiscountPrice != nil {
		setClauses = append(setClauses, fmt.Sprintf("price_discount = $%d", argNum))
		args = append(args, *input.DiscountPrice)
		argNum++
	}

	if input.ExpiresAt != nil {
		setClauses = append(setClauses, fmt.Sprintf("expires_at = $%d", argNum))
		args = append(args, *input.ExpiresAt)
		argNum++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("%s: no fields to update", op)
	}

	args = append(args, boxID)
	whereNum := argNum

	q := database.Query{
		Name: "update_box",
		Sql: fmt.Sprintf(`
		UPDATE business.boxes
		SET %s
		WHERE id = $%d
		RETURNING id, venue_id, category_id, image, price_full, price_discount, created_at, expires_at
	`, strings.Join(setClauses, ", "), whereNum),
	}

	var box entity.Box
	err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&box.ID,
		&box.VenueID,
		&box.CategoryID,
		&box.Image,
		&box.FullPrice,
		&box.DiscountPrice,
		&box.CreatedAt,
		&box.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &box, nil
}

func (r *Repository) GetBoxItems(ctx context.Context, boxID int64, offset, limit int) (pagination.Result[entity.BoxItem], error) {
	const op = "repository.box.GetBoxItems"

	scanner := func(rows pgx.Rows) (entity.BoxItem, error) {
		var item entity.BoxItem

		err := rows.Scan(
			&item.BoxID,
			&item.BoxCode,
			&item.ReservedByUserID,
		)

		if err != nil {
			return entity.BoxItem{}, err
		}

		return item, nil
	}

	p := pagination.Params{
		Table:   "business.box_items",
		Columns: "box_id, box_code, reserved_by_user_id",
		Where:   "box_id = $1",
		Args:    []any{boxID},
		OrderBy: "box_code ASC",
		Offset:  offset,
		Limit:   limit,
	}

	result, err := pagination.List(ctx, r.db.DB(), op, p, scanner)
	if err != nil {
		return pagination.Result[entity.BoxItem]{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
