package mapper

import (
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"
)

func ToCommentResponse(comment entity.Comment) dto.CommentResponse {
	return dto.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		AuthorID:  comment.AuthorID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

func ToCommentWithAuthorResponse(comment entity.CommentWithAuthor, st storage.ObjectStorage) dto.CommentWithAuthorResponse {
	var avatarURL *string
	if comment.AuthorAvatarURL != nil && *comment.AuthorAvatarURL != "" && st != nil {
		url := st.BuildURL(*comment.AuthorAvatarURL)
		avatarURL = &url
	} else {
		avatarURL = comment.AuthorAvatarURL
	}

	return dto.CommentWithAuthorResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		Author: dto.AuthorInfo{
			ID:        comment.AuthorID,
			Username:  comment.AuthorUsername,
			FirstName: comment.AuthorFirstName,
			LastName:  comment.AuthorLastName,
			AvatarURL: avatarURL,
		},
	}
}
