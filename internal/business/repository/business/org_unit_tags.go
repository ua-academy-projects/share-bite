package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
)

func (r *Repository) GetOrgUnitTagSlugs(ctx context.Context, orgUnitID int) ([]string, error) {
	const op = "repository.business.GetOrgUnitTagSlugs"

	q := database.Query{
		Name: "business_repository.GetOrgUnitTagSlugs",
		Sql: `
			SELECT lt.slug
			FROM business.org_unit_tags out
			JOIN business.location_tags lt ON lt.id = out.tag_id
			WHERE out.org_unit_id = $1
			ORDER BY lt.slug
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, orgUnitID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tags = append(tags, slug)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tags, nil
}

func (r *Repository) GetOrgUnitTagsByOrgUnitID(ctx context.Context, ids []int) (map[int][]string, error) {
	const op = "repository.business.GetOrgUnitTagsByOrgUnitID"

	result := make(map[int][]string)
	ids = uniqueInts(ids)
	if len(ids) == 0 {
		return result, nil
	}

	q := database.Query{
		Name: "business_repository.GetOrgUnitTagsByOrgUnitID",
		Sql: `
			SELECT out.org_unit_id, lt.slug
			FROM business.org_unit_tags out
			JOIN business.location_tags lt ON lt.id = out.tag_id
			WHERE out.org_unit_id = ANY($1::int[])
			ORDER BY out.org_unit_id, lt.slug
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var orgUnitID int
		var slug string

		if err := rows.Scan(&orgUnitID, &slug); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		result[orgUnitID] = append(result[orgUnitID], slug)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (r *Repository) SetOrgUnitTagsByIDs(ctx context.Context, orgUnitID int, tagIDs []int) error {
	const op = "repository.business.SetOrgUnitTagsByIDs"

	tagIDs = uniqueInts(tagIDs)

	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	ownTx := !ok
	if ownTx {
		var err error
		tx, err = r.db.DB().BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return fmt.Errorf("%s: begin tx: %w", op, err)
		}
		defer func() {
			if ownTx {
				_ = tx.Rollback(ctx)
			}
		}()
	}

	lockQ := database.Query{
		Name: "business_repository.SetOrgUnitTagsByIDs.LockOrgUnit",
		Sql: `
			SELECT id
			FROM business.org_units
			WHERE id = $1
			FOR UPDATE
		`,
	}

	var lockedID int
	if err := tx.QueryRow(ctx, lockQ.Sql, orgUnitID).Scan(&lockedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return fmt.Errorf("%s: lock org unit: %w", op, err)
	}

	if len(tagIDs) > 0 {
		validateQ := database.Query{
			Name: "business_repository.SetOrgUnitTagsByIDs.Validate",
			Sql: `
				SELECT COUNT(*)
				FROM business.location_tags
				WHERE id = ANY($1::int[])
			`,
		}

		var cnt int
		if err := tx.QueryRow(ctx, validateQ.Sql, tagIDs).Scan(&cnt); err != nil {
			return fmt.Errorf("%s: validate tag ids: %w", op, err)
		}
		if cnt != len(tagIDs) {
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
	}

	deleteQ := database.Query{
		Name: "business_repository.SetOrgUnitTagsByIDs.Delete",
		Sql: `
			DELETE FROM business.org_unit_tags
			WHERE org_unit_id = $1
		`,
	}
	if _, err := tx.Exec(ctx, deleteQ.Sql, orgUnitID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(tagIDs) > 0 {
		insertQ := database.Query{
			Name: "business_repository.SetOrgUnitTagsByIDs.Insert",
			Sql: `
				INSERT INTO business.org_unit_tags (org_unit_id, tag_id)
				SELECT $1, lt.id
				FROM business.location_tags lt
				WHERE lt.id = ANY($2::int[])
				ON CONFLICT (org_unit_id, tag_id) DO NOTHING
			`,
		}

		if _, err := tx.Exec(ctx, insertQ.Sql, orgUnitID, tagIDs); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if ownTx {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("%s: commit tx: %w", op, err)
		}
		ownTx = false
	}
	return nil
}

func (r *Repository) ListLocationTags(ctx context.Context) ([]entity.LocationTag, error) {
	const op = "repository.business.ListLocationTags"

	q := database.Query{
		Name: "business_repository.ListLocationTags",
		Sql: `
			SELECT id, name, slug
			FROM business.location_tags
			ORDER BY id
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	tags := make([]entity.LocationTag, 0)
	for rows.Next() {
		var t entity.LocationTag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tags = append(tags, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tags, nil
}

func uniqueInts(ids []int) []int {
	if len(ids) == 0 {
		return nil
	}

	seen := make(map[int]struct{}, len(ids))
	out := make([]int, 0, len(ids))

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}

	return out
}
