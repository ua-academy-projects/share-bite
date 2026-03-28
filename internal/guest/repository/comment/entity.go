package comment

import (
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type Comment struct {
	ID         int64  `db:"id"`
	PostID     int64  `db:"post_id"`
	CustomerID string `db:"customer_id"`
	Text       string `db:"text"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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

type Comments []Comment

func (cs Comments) ToEntities() []entity.Comment {
	res := make([]entity.Comment, 0, len(cs))
	for i := range cs {
		res = append(res, cs[i].ToEntity())
	}
	return res
}
