package follow

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Follow(ctx context.Context, followerID, followedID string) error {
	sql := `
		INSERT INTO guest.customer_follows(
			follower_customer_id,
			followed_customer_id
		)
		VALUES ($1, $2)
		ON CONFLICT (follower_customer_id, followed_customer_id) DO NOTHING;
	`

	q := database.Query{
		Name: "follow_repository.Follow",
		Sql:  sql,
	}

	_, err := r.db.DB().ExecContext(ctx, q, followerID, followedID)
	if err != nil {
		return executeSQLError(err)
	}

	return nil
}

func (r *Repository) Unfollow(
	ctx context.Context,
	followerID, followedID string,
) error {

	sql := `
		DELETE FROM guest.customer_follows
		WHERE follower_customer_id = $1
		  AND followed_customer_id = $2
	`

	q := database.Query{
		Name: "follow_repository.Unfollow",
		Sql:  sql,
	}

	res, err := r.db.DB().ExecContext(ctx, q, followerID, followedID)
	if err != nil {
		return executeSQLError(err)
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrFollowNotFound
	}

	return nil
}

func (r *Repository) IsFollowing(
	ctx context.Context,
	followerID, followedID string,
) (bool, error) {

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

func (r *Repository) GetAllowedMentions(ctx context.Context, customerID string, mentions []string) ([]string, error) {
	if len(mentions) == 0 {
		return nil, nil
	}

	sql := `
		SELECT followed_customer_id
		FROM guest.customer_follows
		WHERE follower_customer_id = $1
		  AND followed_customer_id = ANY($2::uuid[])
	`

	q := database.Query{
		Name: "follow_repository.GetAllowedMentions",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, customerID, mentions)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, scanRowError(err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *Repository) GetFollowers(ctx context.Context, customerID string) ([]entity.Customer, error) {
	sql := `
		SELECT c.id, c.user_id, c.username, c.first_name, c.last_name, c.avatar_object_key, c.bio, c.is_followers_public, c.is_following_public, c.created_at
		FROM guest.customer_follows cf
		JOIN guest.customers c ON cf.follower_customer_id = c.id
		WHERE cf.followed_customer_id = $1
	`
	q := database.Query{
		Name: "follow_repository.GetFollowers",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, customerID)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var resultRows []FollowerCustomer
	if err := pgxscan.ScanAll(&resultRows, rows); err != nil {
		return nil, scanRowsError(err)
	}

	result := make([]entity.Customer, 0, len(resultRows))
	for _, row := range resultRows {
		result = append(result, row.ToEntity())
	}

	return result, nil
}
