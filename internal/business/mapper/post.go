package mapper

import (
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func ToPostResponse(post *entity.Post) dto.PostResponse {
	return dto.PostResponse{
		ID:       post.ID,
		Content:  post.Content,
		ImageURL: post.ImageURL,
	}
}
