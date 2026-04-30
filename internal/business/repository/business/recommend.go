package business

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

func (r *Repository) GetTopTagsByUserLikes(ctx context.Context, userID string, tagsToFetch int) ([]string, error) {
	const op = "repository.business.GetTopTagsByUserLikes"

	// Note: Adjust table names `guest.likes` and `business.venue_tags` to match your exact schema.
	sql := `
		SELECT vt.tag
		FROM guest.post_likes l
		JOIN business.posts p ON l.post_id = p.id
		JOIN business.venue_tags vt ON p.org_id = vt.venue_id
		WHERE l.customer_id = $1::uuid
		GROUP BY vt.tag
		ORDER BY COUNT(vt.tag) DESC
		LIMIT $2
	`

	q := database.Query{
		Name: "business_repository.GetTopTagsByUserLikes",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, userID, tagsToFetch)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w", op, err)
	}

	return tags, nil
}

func (r *Repository) GetVenuesByTag(ctx context.Context, tag string, quota int, seenIDs []int) ([]entity.OrgUnit, error) {
	const op = "repository.business.GetVenuesByTag"

	sql := `
		SELECT ou.id, ou.org_account_id, ou.profile_type, ou.name, ou.avatar, ou.banner, 
		       ou.description, ou.parent_id, ou.latitude, ou.longitude, ou.h3_hash
		FROM business.org_units ou
		JOIN business.venue_tags vt ON ou.id = vt.venue_id
		WHERE vt.tag = $1 
		  AND ou.profile_type = 'VENUE'
		  AND ou.id <> ALL($2::int[])
		LIMIT $3
	`

	q := database.Query{
		Name: "business_repository.GetVenuesByTag",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, tag, seenIDs, quota)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanOrgUnits(rows, op)
}

func (r *Repository) GetRandomVenues(ctx context.Context, deficit int, seenIDs []int) ([]entity.OrgUnit, error) {
	const op = "repository.business.GetRandomVenues"

	sql := `
		SELECT id, org_account_id, profile_type, name, avatar, banner, 
		       description, parent_id, latitude, longitude, h3_hash
		FROM business.org_units
		WHERE profile_type = 'VENUE' 
		  AND id <> ALL($1::int[])
		ORDER BY RANDOM()
		LIMIT $2
	`

	q := database.Query{
		Name: "business_repository.GetRandomVenues",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, seenIDs, deficit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanOrgUnits(rows, op)
}

func scanOrgUnits(rows pgx.Rows, op string) ([]entity.OrgUnit, error) {
	var result []entity.OrgUnit

	for rows.Next() {
		var ou OrgUnit
		err := rows.Scan(
			&ou.Id, &ou.OrgAccountId, &ou.ProfileType, &ou.Name, &ou.Avatar, &ou.Banner,
			&ou.Description, &ou.ParentId, &ou.Latitude, &ou.Longitude, &ou.H3Hash,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		result = append(result, ou.ToEntity())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w", op, err)
	}

	return result, nil
}
