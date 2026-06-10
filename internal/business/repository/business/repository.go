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
	constraintOrgAccountIDKey = "org_units_brand_owner_uidx"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrNoAvailableItems = errors.New("no available box items")
	ErrInvalidStatus    = errors.New("organization unit is not in a resubmittable state")
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
		name,
		description,
		avatar,
		banner
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	q := database.Query{
		Name: "business_repository.Create",
		Sql:  sql,
	}

	args := []any{
		in.OrgAccountId,
		in.ProfileType,
		in.Name,
		in.Description,
		in.Avatar,
		in.Banner,
	}

	var businessOrgID int
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&businessOrgID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				switch pgErr.ConstraintName {
				case constraintOrgAccountIDKey:
					return 0, apperror.BadRequest("org account already has a business unit")
				default:
					return 0, fmt.Errorf("unknown unique rule violated: %s", pgErr.ConstraintName)
				}
			}
		}
		return 0, scanRowError(err)
	}

	return businessOrgID, nil
}

func (r *Repository) UpdateOrg(ctx context.Context, id int, orgAccountID uuid.UUID, in entity.UpdateOrgUnitInput) (*entity.OrgUnit, error) {
	sql := `
		UPDATE business.org_units
		SET
			name = COALESCE($1, name),
			avatar = COALESCE($2, avatar),
			banner = COALESCE($3, banner),
			description = COALESCE($4, description),
		    status = 'pending'
		WHERE id = $5 AND org_account_id = $6 AND profile_type = 'BRAND'
		RETURNING id, org_account_id, profile_type, name, avatar, banner, description, status
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
		&ou.Status,
	)
	if err != nil {
		if errors.Is(err, stdsql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
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
		WHERE id = $1 AND org_account_id = $2 AND profile_type = 'BRAND'
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
