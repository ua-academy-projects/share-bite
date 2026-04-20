package business

import (
	"context"
	stdsql "database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

const (
	constraintOrgAccountIDKey = "org_units_org_account_id_key"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, in entity.OrgUnit) (int, error) {
	sql := `
	INSERT INTO business.org_units(
		org_account_id,
		profile_type,
		parent_id,
		name,
		description,
		latitude,
		longitude
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	q := database.Query{
		Name: "business_repository.Create",
		Sql:  sql,
	}

	args := []any{
		in.OrgAccountId,
		in.ProfileType,
		in.ParentId,
		in.Name,
		in.Description,
		in.Latitude,
		in.Longitude,
	}

	var businessOrgID int
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&businessOrgID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case constraintOrgAccountIDKey:
					return 0, apperror.ErrOrgAccountAlreadyHasUnit
				default:
					return 0, fmt.Errorf("Unknown unique rule violated: %s", pgErr.ConstraintName)
				}

			}
		}
		return 0, scanRowError(err)
	}

	return businessOrgID, nil
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

func (r *Repository) UpdateOrg(ctx context.Context, id int, orgAccountID uuid.UUID, in entity.UpdateOrgUnitInput) (*entity.OrgUnit, error) {
	sql := `
		UPDATE business.org_units
		SET
			name = COALESCE($1, name),
			avatar = COALESCE($2, avatar),
			banner = COALESCE($3, banner),
			description = COALESCE($4, description),
			latitude = COALESCE($5, latitude),
			longitude = COALESCE($6, longitude)
		WHERE id = $7 AND org_account_id = $8
		RETURNING id, org_account_id, profile_type, name, avatar, banner, description, parent_id, latitude, longitude
	`
	q := database.Query{
		Name: "business_repository.UpdateOrg",
		Sql:  sql,
	}

	row := r.db.DB().QueryRowContext(
		ctx,
		q,
		in.Name,
		in.Avatar,
		in.Banner,
		in.Description,
		in.Latitude,
		in.Longitude,
		id,
		orgAccountID,
	)

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
		if errors.Is(err, stdsql.ErrNoRows) {
			return nil, ErrNotFound
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, scanRowError(err)
	}

	result := ou.ToEntity()
	return &result, nil
}

func (r *Repository) DeleteOrg(ctx context.Context, id int, orgAccountID uuid.UUID) error {
	sql := `
		DELETE FROM business.org_units
		WHERE id = $1 AND org_account_id = $2
	`
	q := database.Query{
		Name: "business_repository.DeleteOrg",
		Sql:  sql,
	}

	rows, err := r.db.DB().ExecContext(ctx, q, id, orgAccountID)
	if err != nil {
		return executeSQLError(err)
	}

	if rows.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}


func (r *Repository) GetOrgIDByUserID(ctx context.Context, userID int64) (int, error) {
	q := database.Query{
		Name: "get_org_by_user_id",
		Sql: `
			SELECT id
			FROM business.org_units
			WHERE org_account_id = $1
		`,
	}

	var orgID int

	err := r.db.DB().QueryRowContext(ctx, q, userID).Scan(&orgID)
	if err != nil {
		if err == stdsql.ErrNoRows {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return orgID, nil
}


}
