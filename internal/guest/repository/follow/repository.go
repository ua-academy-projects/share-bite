package follow

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Follow(ctx context.Context, followerID, followedID string) (entity.CustomerFollow, error) {
	sql := `
		INSERT INTO guest.customer_follows(
			follower_customer_id,
			followed_customer_id
		) VALUES ($1, $2)
		RETURNING id, follower_customer_id, followed_customer_id, created_at
	`

	q := database.Query{
		Name: "follow_repository.Follow",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, followerID, followedID)
	if err != nil {
		return entity.CustomerFollow{}, executeSQLError(err)
	}
	defer rows.Close()

	var follow CustomerFollow
	if err := pgxscan.ScanOne(&follow, rows); err != nil {
		return entity.CustomerFollow{}, scanRowError(err)
	}

	return follow.ToEntity(), nil
}
func (r *Repository) Unfollow(ctx context.Context, followerID, followedID string) error {
	sql := `
		DELETE FROM guest.customer_follows
		WHERE follower_customer_id = $1
		  AND followed_customer_id = $2
	`

	q := database.Query{
		Name: "follow_repository.Unfollow",
		Sql:  sql,
	}

	_, err := r.db.DB().ExecContext(ctx, q, followerID, followedID)
	if err != nil {
		return executeSQLError(err)
	}

	return nil
}
func (r *Repository) IsFollowing(ctx context.Context, followerID, followedID string) (bool, error) {
	sql := `
		SELECT EXISTS (
			SELECT 1
			FROM guest.customer_follows
			WHERE follower_customer_id = $1
			  AND followed_customer_id = $2
		)
	`

	q := database.Query{
		Name: "follow_repository.IsFollowing",
		Sql:  sql,
	}

	var exists bool
	if err := r.db.DB().QueryRowContext(ctx, q, followerID, followedID).Scan(&exists); err != nil {
		return false, scanRowError(err)
	}

	return exists, nil
}
