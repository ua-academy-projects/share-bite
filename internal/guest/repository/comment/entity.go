package comment

import (
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type Comment struct {
	ID         int64     `db:"id"`
	PostID     int64     `db:"post_id"`
	CustomerID string    `db:"customer_id"`
	Text       string    `db:"text"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (c *Comment) ToEntity() entity.Comment {
	return entity.Comment{
		ID:         c.ID,
		PostID:     c.PostID,
		CustomerID: c.CustomerID,
		Text:       c.Text,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

type CommentRow struct {
	ID         int64     `db:"id"`
	PostID     int64     `db:"post_id"`
	CustomerID string    `db:"customer_id"`
	Text       string    `db:"text"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`

	CustID        string  `db:"cust_id"`
	CustUserID    string  `db:"cust_user_id"`
	CustUserName  string  `db:"cust_username"`
	CustFirstName string  `db:"cust_first_name"`
	CustLastName  string  `db:"cust_last_name"`
	CustAvatar    *string `db:"cust_avatar"`
	CustBio       *string `db:"cust_bio"`
}

func (r CommentRow) ToEntity() entity.CommentWithCustomer {
	return entity.CommentWithCustomer{
		Comment: entity.Comment{
			ID:         r.ID,
			PostID:     r.PostID,
			CustomerID: r.CustomerID,
			Text:       r.Text,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		},
		Customer: entity.Customer{
			ID:              r.CustID,
			UserID:          r.CustUserID,
			UserName:        r.CustUserName,
			FirstName:       r.CustFirstName,
			LastName:        r.CustLastName,
			AvatarObjectKey: r.CustAvatar,
			Bio:             r.CustBio,
		},
	}
}
