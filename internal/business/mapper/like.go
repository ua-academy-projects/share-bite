package mapper

import (
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func ToLikeItem(like entity.LikeWithAuthor) dto.LikeItem {
	return dto.LikeItem{
		ID:              like.ID,
		PostID:          like.PostID,
		AuthorID:        like.AuthorID,
		AuthorUsername:  like.AuthorUsername,
		AuthorFirstName: like.AuthorFirstName,
		AuthorLastName:  like.AuthorLastName,
		AuthorAvatarURL: like.AuthorAvatarURL,
		CreatedAt:       like.CreatedAt,
	}
}
