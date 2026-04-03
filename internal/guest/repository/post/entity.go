package post

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type Post struct {
	ID         int64  `db:"id"`
	CustomerID string `db:"customer_id"`
	VenueID    string `db:"venue_id"`
	Text       string `db:"text"`
	Rating     int16  `db:"rating"`
	Status     string `db:"status"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (p *Post) ToEntity() entity.Post {
	return entity.Post{
		ID:         strconv.FormatInt(p.ID, 10),
		CustomerID: p.CustomerID,
		VenueID:    p.VenueID,
		Text:       p.Text,
		Rating:     p.Rating,
		Status:     p.Status,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

type Posts []Post

func (ps Posts) ToEntities() []entity.Post {
	res := make([]entity.Post, 0, len(ps))
	for i := range ps {
		res = append(res, ps[i].ToEntity())
	}
	return res
}

func executeSQLError(err error) error {
	return fmt.Errorf("execute sql: %w", err)
}

func scanRowError(err error) error {
	return fmt.Errorf("scan row: %w", err)
}

func scanRowsError(err error) error {
	return fmt.Errorf("scan rows: %w", err)
}
