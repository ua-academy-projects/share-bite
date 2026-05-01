package guest

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/ua-academy-projects/share-bite/internal/admin-auth/dto"
	apperr "github.com/ua-academy-projects/share-bite/internal/admin-auth/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Client struct {
	client database.Client
}

func NewClient(client database.Client) *Client {
	return &Client{client: client}
}

func (c *Client) GetCustomerByUserID(ctx context.Context, userID string) (*dto.CustomerProfileData, error) {
	q := database.Query{
		Name: "guest.GetCustomerByUserID",
		Sql: `
			SELECT username, first_name, last_name, avatar_object_key, bio 
			FROM guest.customers 
			WHERE user_id = $1
		`,
	}

	row := c.client.DB().QueryRowContext(ctx, q, userID)

	p := new(dto.CustomerProfileData)
	if err := row.Scan(
		&p.Username,
		&p.FirstName,
		&p.LastName,
		&p.AvatarKey,
		&p.Bio,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query customer profile", err)
	}

	return p, nil
}
