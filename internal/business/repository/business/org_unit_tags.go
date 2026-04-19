package business

import (
	"context"
	"fmt"
	"strings"

	"github.com/ua-academy-projects/share-bite/pkg/database"
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

func (r *Repository) SetOrgUnitTagsBySlugs(ctx context.Context, orgUnitID int, slugs []string) error {
	const op = "repository.business.SetOrgUnitTagsBySlugs"

	slugs = normalizeAndUniqueSlugs(slugs)

	if len(slugs) > 0 {
		validateQ := database.Query{
			Name: "business_repository.SetOrgUnitTagsBySlugs.Validate",
			Sql: `
				SELECT COUNT(*)
				FROM business.location_tags
				WHERE slug = ANY($1::text[])
			`,
		}

		var cnt int
		if err := r.db.DB().QueryRowContext(ctx, validateQ, slugs).Scan(&cnt); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if cnt != len(slugs) {
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}
	}

	deleteQ := database.Query{
		Name: "business_repository.SetOrgUnitTagsBySlugs.Delete",
		Sql: `
			DELETE FROM business.org_unit_tags
			WHERE org_unit_id = $1
		`,
	}

	if _, err := r.db.DB().ExecContext(ctx, deleteQ, orgUnitID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(slugs) == 0 {
		return nil
	}

	insertQ := database.Query{
		Name: "business_repository.SetOrgUnitTagsBySlugs.Insert",
		Sql: `
			INSERT INTO business.org_unit_tags (org_unit_id, tag_id)
			SELECT $1, lt.id
			FROM business.location_tags lt
			WHERE lt.slug = ANY($2::text[])
			ON CONFLICT (org_unit_id, tag_id) DO NOTHING
		`,
	}

	if _, err := r.db.DB().ExecContext(ctx, insertQ, orgUnitID, slugs); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func normalizeAndUniqueSlugs(slugs []string) []string {
	if len(slugs) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(slugs))
	out := make([]string, 0, len(slugs))

	for _, slug := range slugs {
		s := strings.TrimSpace(strings.ToLower(slug))
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}

	return out
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
