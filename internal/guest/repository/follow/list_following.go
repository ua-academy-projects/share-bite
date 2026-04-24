package follow

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"strconv"
	"time"
)

func (r *Repository) ListFollowing(
	ctx context.Context,
	customerID string,
	cursorTime time.Time,
	cursorID string,
	limit int,
) ([]entity.Follower, error) {

	sql := `
		SELECT
			cf.id AS follow_id,
			cf.created_at AS follow_created_at,
			c.id,
			c.user_id,
			c.username,
			c.first_name,
			c.last_name,
			c.avatar_object_key,
			c.bio,
			c.is_followers_public,
			c.is_following_public,
			c.created_at
		FROM guest.customer_follows cf
		JOIN guest.customers c
			ON c.id = cf.followed_customer_id
		WHERE cf.follower_customer_id = $1
	`

	args := []any{customerID}

	if !cursorTime.IsZero() && cursorID != "" {
		sql += `
			AND (cf.created_at, cf.id) < ($2, $3)
		`
		args = append(args, cursorTime, cursorID)
	}

	sql += `
		ORDER BY cf.created_at DESC, cf.id DESC
		LIMIT $` + strconv.Itoa(len(args)+1)

	args = append(args, limit)

	q := database.Query{
		Name: "follow_repository.ListFollowing",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var resultRows []FollowCustomerRow
	if err := pgxscan.ScanAll(&resultRows, rows); err != nil {
		return nil, scanRowsError(err)
	}

	following := make([]entity.Follower, 0, len(resultRows))
	for _, row := range resultRows {
		following = append(following, row.ToEntity())
	}

	return following, nil
}
