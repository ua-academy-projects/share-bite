package mapper

import (
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func ToLikeItem(like entity.Like) dto.LikeItem {
	return dto.LikeItem{
		ID:         like.ID,
		PostID:     like.PostID,
		CustomerID: like.CustomerID,
		CreatedAt:  like.CreatedAt,
	}
}
