package mapper

import (
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"
)

func ToLikeItem(like entity.LikeWithAuthor, st storage.ObjectStorage) dto.LikeItem {
	var avatarURL *string
	if like.AuthorAvatarURL != nil && *like.AuthorAvatarURL != "" && st != nil {
		url := st.BuildURL(*like.AuthorAvatarURL)
		avatarURL = &url
	} else {
		avatarURL = like.AuthorAvatarURL
	}

	return dto.LikeItem{
		ID:              like.ID,
		PostID:          like.PostID,
		AuthorID:        like.AuthorID,
		AuthorUsername:  like.AuthorUsername,
		AuthorFirstName: like.AuthorFirstName,
		AuthorLastName:  like.AuthorLastName,
		AuthorAvatarURL: avatarURL,
		CreatedAt:       like.CreatedAt,
	}
}
