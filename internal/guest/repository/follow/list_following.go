package follow

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"strconv"
	"time"
)

func (r *Repository) ListFollowingEnriched(ctx context.Context, requesterID string, customerID string, cursorTime time.Time, cursorID string, limit int) ([]entity.Follower, error) {
	sql := `
		SELECT
			cf.id AS follow_id,
			cf.created_at AS follow_created_at,
		
			c.id,
			c.avatar_object_key,
			c.is_followers_public,
			c.is_following_public,
		
			(f1.follower_customer_id IS NOT NULL) AS is_following,
			(f2.follower_customer_id IS NOT NULL) AS is_followed_by,
			(f1.follower_customer_id IS NOT NULL AND f2.follower_customer_id IS NOT NULL) AS is_mutual
		
		FROM guest.customer_follows cf
		
		JOIN guest.customers c
			ON c.id = cf.followed_customer_id
		
		-- requester -> target
		LEFT JOIN guest.customer_follows f1
			ON f1.followed_customer_id = c.id
		   AND f1.follower_customer_id = $1
		
		-- target -> requester
		LEFT JOIN guest.customer_follows f2
			ON f2.followed_customer_id = $1
		   AND f2.follower_customer_id = c.id
		
		WHERE cf.follower_customer_id = $2
		`

	args := []any{requesterID, customerID}

	if !cursorTime.IsZero() && cursorID != "" {
		sql += `AND (cf.created_at, cf.id) < ($3, $4)`
		args = append(args, cursorTime, cursorID)
	}

	sql += `ORDER BY cf.created_at DESC, cf.id DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	q := database.Query{
		Name: "follow_repository.ListFollowingEnriched",
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
