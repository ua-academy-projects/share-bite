package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

func (r *Repository) GetById(ctx context.Context, id int) (*entity.OrgUnit, error) {
	const op = "repository.business.GetById"
	sql := `
		SELECT id, org_account_id, profile_type, name, avatar, banner, description, parent_id, latitude, longitude
		FROM business.org_units
		WHERE id = $1
	`
	q := database.Query{
		Name: "business_repository.GetById",
		Sql:  sql,
	}

	row := r.db.DB().QueryRowContext(ctx, q, id)

	var ou OrgUnit
	err := row.Scan(
		&ou.Id,
		&ou.OrgAccountId,
		&ou.ProfileType,
		&ou.Name,
		&ou.Avatar,
		&ou.Banner,
		&ou.Description,
		&ou.ParentId,
		&ou.Latitude,
		&ou.Longitude,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := ou.ToEntity()
	return &result, nil
}

func (r *Repository) ListByParentID(ctx context.Context, parentID, offset, limit int) (pagination.Result[entity.OrgUnit], error) {
	const op = "repository.business.ListByParentID"
	units, err := pagination.List(ctx, r.db.DB(), "business_repository.ListByParentID",
		pagination.Params{
			Table:   "business.org_units",
			Columns: "id, org_account_id, profile_type, name, avatar, banner, description, parent_id, latitude, longitude",
			Where:   "parent_id = $1",
			OrderBy: "id",
			Args:    []any{parentID},
			Offset:  offset,
			Limit:   limit,
		},
		func(rows pgx.Rows) (entity.OrgUnit, error) {
			var ou OrgUnit
			err := rows.Scan(
				&ou.Id,
				&ou.OrgAccountId,
				&ou.ProfileType,
				&ou.Name,
				&ou.Avatar,
				&ou.Banner,
				&ou.Description,
				&ou.ParentId,
				&ou.Latitude,
				&ou.Longitude,
			)
			if err != nil {
				return entity.OrgUnit{}, fmt.Errorf("%s: %w", op, err)
			}
			return ou.ToEntity(), nil
		},
	)
	if err != nil {
		return pagination.Result[entity.OrgUnit]{}, fmt.Errorf("%s: %w", op, err)
	}

	return units, nil
}

func (r *Repository) GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error) {
	const op = "repository.business.GetVenuesByIDs"
	q := database.Query{
		Name: "business_repository.GetVenuesByIDs",
		Sql: `
			SELECT id, name, description, avatar, banner
			FROM business.org_units
			WHERE id = ANY($1) AND parent_id IS NOT NULL
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []entity.OrgUnit
	for rows.Next() {
		var ou OrgUnit
		err := rows.Scan(
			&ou.Id,
			&ou.Name,
			&ou.Description,
			&ou.Avatar,
			&ou.Banner,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, ou.ToEntity())
	}

	return result, nil
}
