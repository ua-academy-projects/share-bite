package business

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

func (c *Client) GetBusinessByUserID(ctx context.Context, userID string) (*dto.BusinessProfileData, error) {
	q := database.Query{
		Name: "business.GetBusinessByUserID",
		Sql: `
          SELECT 
              profile_type, 
              name, 
              COALESCE(avatar, ''), 
              COALESCE(banner, ''), 
              COALESCE(description, ''), 
              latitude,
              longitude
          FROM business.org_units 
          WHERE org_account_id = $1
          ORDER BY CASE WHEN profile_type = 'BRAND' THEN 1 ELSE 2 END ASC
          LIMIT 1
       `,
	}

	row := c.client.DB().QueryRowContext(ctx, q, userID)

	p := new(dto.BusinessProfileData)
	if err := row.Scan(
		&p.ProfileType,
		&p.Name,
		&p.Avatar,
		&p.Banner,
		&p.Description,
		&p.Latitude,
		&p.Longitude,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query business profile", err)
	}

	return p, nil
}
