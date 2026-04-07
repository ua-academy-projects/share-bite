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
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("%w", err)
	}

	result := ou.ToEntity()
	return &result, nil
}

func (r *Repository) ListByParentID(ctx context.Context, parentID, offset, limit int) (pagination.Result[entity.OrgUnit], error) {
	return pagination.List(ctx, r.db.DB(), "business_repository.ListByParentID",
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
				return entity.OrgUnit{}, err
			}
			return ou.ToEntity(), nil
		},
	)
}
