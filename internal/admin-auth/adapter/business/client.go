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
              longitude,
              status,
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
		&p.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, apperr.Wrap(http.StatusInternalServerError, "failed to query business profile", err)
	}

	return p, nil
}

func (c *Client) GetPendingBusinesses(ctx context.Context, limit, offset int) ([]dto.PendingBusinessListItem, int, error) {
	q := database.Query{
		Name: "business.GetPendingBusinesses",
		Sql: `
          SELECT 
              id,
              org_account_id,
              name,
              COALESCE(avatar, ''),
              COALESCE(description, ''),
              status,
              COUNT(*) OVER() AS total_count
          FROM business.org_units
          WHERE status = 'pending'
          ORDER BY id ASC
          LIMIT $1 OFFSET $2
       `,
	}

	rows, err := c.client.DB().QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, apperr.Wrap(http.StatusInternalServerError, "failed to get pending businesses", err)
	}
	defer rows.Close()

	var items []dto.PendingBusinessListItem
	var totalCount int

	for rows.Next() {
		var item dto.PendingBusinessListItem
		if err := rows.Scan(
			&item.ID, &item.OrgAccountID, &item.Name,
			&item.Avatar, &item.Description, &item.Status, &totalCount,
		); err != nil {
			return nil, 0, apperr.Wrap(http.StatusInternalServerError, "failed to scan pending business", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperr.Wrap(http.StatusInternalServerError, "rows iteration error", err)
	}

	if items == nil {
		items = make([]dto.PendingBusinessListItem, 0)
	}

	return items, totalCount, nil
}

func (c *Client) ReviewBusiness(ctx context.Context, params dto.ReviewBusinessParams) error {
	q := database.Query{
		Name: "business.ReviewBusiness",
		Sql: `
          WITH prev AS (
             SELECT id, status FROM business.org_units WHERE id = $1
          ),
          upd AS (
             UPDATE business.org_units
             SET status = $2
             WHERE id = $1
             RETURNING id
          )
          INSERT INTO business.verification_logs (org_unit_id, admin_id, old_status, new_status, comment)
          SELECT p.id, $3, p.status, $2, $4
          FROM prev p;
       `,
	}

	result, err := c.client.DB().ExecContext(
		ctx,
		q,
		params.OrgUnitID,
		params.NewStatus,
		params.AdminID,
		params.Comment,
	)
	if err != nil {
		return apperr.Wrap(http.StatusInternalServerError, "failed to review business and create log", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("business org_unit not found")
	}

	return nil
}

func (c *Client) GetBusinessStatusAndOwner(ctx context.Context, orgUnitID int) (string, string, error) {
	q := database.Query{
		Name: "business.GetBusinessStatus",
		Sql:  `SELECT status, org_account_id FROM business.org_units WHERE id = $1`,
	}
	var status, ownerID string
	err := c.client.DB().QueryRowContext(ctx, q, orgUnitID).Scan(&status, &ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", apperr.ErrBusinessNotFound
		}
		return "", "", apperr.Wrap(http.StatusInternalServerError, "failed to get business status from database", err)
	}
	return status, ownerID, nil
}
