package business

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

var ErrNotFound = errors.New("not found")

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
			return nil, ErrNotFound
		}

		return nil, scanRowError(err)
	}

	result := ou.ToEntity()
	return &result, nil
}

func (r *Repository) ListByParentID(ctx context.Context, parentID, offset, limit int) ([]entity.OrgUnit, error) {
	sql := `
		SELECT id, org_account_id, profile_type, name, avatar, banner, description, parent_id, latitude, longitude
		FROM business.org_units
		WHERE parent_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`
	q := database.Query{
		Name: "business_repository.ListByParentID",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, parentID, limit, offset)
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
