package follow

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	customerrepo "github.com/ua-academy-projects/share-bite/internal/guest/repository/customer"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

func (r *Repository) ListFollowing(ctx context.Context, customerID string) ([]entity.Customer, error) {
	sql := `
		SELECT
			c.id,
			c.user_id,
			c.username,
			c.first_name,
			c.last_name,
			c.avatar_object_key,
			c.bio,
			c.created_at
		FROM guest.customer_follows cf
		JOIN guest.customers c
			ON c.id = cf.followed_customer_id
		WHERE cf.follower_customer_id = $1
		ORDER BY c.created_at DESC
	`

	q := database.Query{
		Name: "follow_repository.ListFollowing",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, customerID)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var customers customerrepo.Customers
	if err := pgxscan.ScanAll(&customers, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return customers.ToEntities(), nil
}
