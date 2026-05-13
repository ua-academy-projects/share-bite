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
	sql := `
		WITH all_user_liked_orgs AS (
			SELECT gp.venue_id::int AS org_id
			FROM guest.post_likes gl
			JOIN guest.posts gp ON gl.post_id = gp.id
			WHERE gl.customer_id = $1::uuid

			UNION ALL

			SELECT bp.org_id
			FROM business.likes bl
			JOIN business.posts bp ON bl.post_id = bp.id
			WHERE bl.author_id = $1::uuid
		)
		SELECT lt.name
		FROM all_user_liked_orgs alo
		JOIN business.org_unit_tags out ON alo.org_id = out.org_unit_id
		JOIN business.location_tags lt ON out.tag_id = lt.id
		GROUP BY lt.name
		ORDER BY COUNT(lt.name) DESC
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

func (r *Repository) GetPostsByTag(ctx context.Context, tag string, quota int, seenCompositeIDs []string, h3Hashes []string) ([]entity.RecomendedPost, error) {
	const op = "repository.business.GetPostsByTag"

	sql := `
		WITH unified_posts AS (
			SELECT id, venue_id::int AS org_id, text AS content, 'guest' AS post_type, created_at
			FROM guest.posts
			WHERE status = 'published'
			
			UNION ALL
			
			SELECT id, org_id, content, 'business' AS post_type, created_at
			FROM business.posts
		)
		SELECT up.id, up.org_id, up.content, up.post_type, up.created_at
		FROM unified_posts up
		JOIN business.org_units ou ON up.org_id = ou.id
		JOIN business.org_unit_tags out ON ou.id = out.org_unit_id
		JOIN business.location_tags lt ON out.tag_id = lt.id
		WHERE lt.name = $1
		  AND NOT (up.post_type || ':' || up.id::text) = ANY($2::text[]) 
		  AND ou.h3_hash = ANY($4::text[])
		ORDER BY up.created_at DESC
		LIMIT $3
	`
	q := database.Query{
		Name: "business_repository.GetPostsByTag",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, tag, seenCompositeIDs, quota, h3Hashes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanPosts(rows, op)
}

func (r *Repository) GetRandomPosts(ctx context.Context, deficit int, seenCompositeIDs []string, h3Hashes []string) ([]entity.RecomendedPost, error) {
	const op = "repository.business.GetRandomPosts"
	sql := `
		WITH unified_posts AS (
			SELECT id, venue_id::int AS org_id, text AS content, 'guest' AS post_type, created_at
			FROM guest.posts
			WHERE status = 'published'
			
			UNION ALL
			
			SELECT id, org_id, content, 'business' AS post_type, created_at
			FROM business.posts
		)
		SELECT up.id, up.org_id, up.content, up.post_type, up.created_at
		FROM unified_posts up
		JOIN business.org_units ou ON up.org_id = ou.id
		WHERE NOT (up.post_type || ':' || up.id::text) = ANY($1::text[])
		  AND ou.h3_hash = ANY($3::text[])
		ORDER BY RANDOM()
		LIMIT $2
	`
	q := database.Query{
		Name: "business_repository.GetRandomPosts",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, seenCompositeIDs, deficit, h3Hashes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanPosts(rows, op)
}

func scanPosts(rows pgx.Rows, op string) ([]entity.RecomendedPost, error) {
	var result []entity.RecomendedPost

	for rows.Next() {
		var p entity.RecomendedPost
		err := rows.Scan(
			&p.ID, &p.OrgID, &p.Content, &p.PostType, &p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows err: %w", op, err)
	}

	return result, nil
}
