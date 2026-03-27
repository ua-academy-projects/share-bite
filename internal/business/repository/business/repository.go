package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{
		db: db,
	}
}

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
			return nil, fmt.Errorf("org unit with id %d was not found", id)
		}

		return nil, scanRowError(err)
	}

	result := ou.ToEntity()
	return &result, nil
}

func (r *Repository) List(ctx context.Context, offset, limit int) ([]entity.OrgUnit, error) {
	sql := `
		SELECT id, org_account_id, profile_type, name, avatar, banner, description, parent_id, latitude, longitude
		FROM business.org_units
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	q := database.Query{
		Name: "business_repository.List",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var orgUnits []entity.OrgUnit
	for rows.Next() {
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
			return nil, scanRowsError(err)
		}

		orgUnits = append(orgUnits, ou.ToEntity())
	}

	if err := rows.Err(); err != nil {
		return nil, scanRowsError(err)
	}

	return orgUnits, nil
}
