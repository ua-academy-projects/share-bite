package business

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
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

		return nil, scanRowError(err)
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

func (r *Repository) GetBrandIDByOwnerUserID(ctx context.Context, userID string) (int, error) {
	sql := `
		SELECT id
		FROM business.org_units
		WHERE org_account_id = $1::uuid
		  AND profile_type = 'BRAND'
		LIMIT 1
	`
	q := database.Query{
		Name: "business_repository.GetBrandIDByOwnerUserID",
		Sql:  sql,
	}

	var brandID int
	err := r.db.DB().QueryRowContext(ctx, q, userID).Scan(&brandID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, scanRowError(err)
	}

	return brandID, nil
}

func (r *Repository) CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error) {
	sql := `
		INSERT INTO business.org_units (
			org_account_id, profile_type, parent_id, name, avatar, banner, description, latitude, longitude
		)
		VALUES ($1::uuid, 'VENUE', $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, org_account_id, profile_type, name, avatar, banner, description, parent_id, latitude, longitude
	`
	q := database.Query{
		Name: "business_repository.CreateLocation",
		Sql:  sql,
	}

	var ou OrgUnit
	err := r.db.DB().QueryRowContext(
		ctx, q,
		ownerUserID,
		brandID,
		in.Name,
		in.Avatar,
		in.Banner,
		in.Description,
		in.Latitude,
		in.Longitude,
	).Scan(
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
		return nil, scanRowError(err)
	}

	result := ou.ToEntity()
	return &result, nil
}

func (r *Repository) UpdateLocation(ctx context.Context, locationID int, brandID int, in dto.UpdateLocationInput) (*entity.OrgUnit, error) {
	sql := `
		UPDATE business.org_units
		SET
			name = COALESCE($1, name),
			avatar = COALESCE($2, avatar),
			banner = COALESCE($3, banner),
			description = COALESCE($4, description),
			latitude = COALESCE($5, latitude),
			longitude = COALESCE($6, longitude)
		WHERE id = $7
		  AND parent_id = $8
		  AND profile_type = 'VENUE'
		RETURNING id, org_account_id, profile_type, name, avatar, banner, description, parent_id, latitude, longitude
	`
	q := database.Query{
		Name: "business_repository.UpdateLocation",
		Sql:  sql,
	}

	var ou OrgUnit
	err := r.db.DB().QueryRowContext(
		ctx, q,
		in.Name,
		in.Avatar,
		in.Banner,
		in.Description,
		in.Latitude,
		in.Longitude,
		locationID,
		brandID,
	).Scan(
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
		return nil, scanRowError(err)
	}

	result := ou.ToEntity()
	return &result, nil
}

func (r *Repository) DeleteLocation(ctx context.Context, locationID int, brandID int) error {
	sql := `
		DELETE FROM business.org_units
		WHERE id = $1
		  AND parent_id = $2
		  AND profile_type = 'VENUE'
	`
	q := database.Query{
		Name: "business_repository.DeleteLocation",
		Sql:  sql,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, locationID, brandID)
	if err != nil {
		return executeSQLError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
