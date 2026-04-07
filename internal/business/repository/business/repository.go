package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"



	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"

	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
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

func (r *Repository) CheckOwnership(ctx context.Context, userID string, unitID int) error {
	checkQuery := `SELECT id
					FROM business.org_units
					WHERE id = $1
					  AND (
					  	org_account_id = $2
						OR 
						parent_id IN (
							SELECT id FROM business.org_units WHERE org_account_id = $2
						)
					);`

	q := database.Query{
		Name: "check_ownership.CheckOwnership",
		Sql:  checkQuery,
	}

	var foundID int
	var err error

	if tx, ok := ctx.Value(pg.TxKey).(pgx.Tx); ok {
		err = tx.QueryRow(ctx, checkQuery, unitID, userID).Scan(&foundID)
	} else {
		err = r.db.DB().QueryRowContext(ctx, q, unitID, userID).Scan(&foundID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return biserr.ErrForbidden
		}
		return fmt.Errorf("execute check ownership query: %w", err)
	}
	return nil
}
