package mapper

import (
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func ToPostResponse(post *entity.PostWithPhotos) dto.PostResponse {
	return dto.PostResponse{
		ID:        post.ID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Images:    post.Images,

		Org: struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			ProfileType string `json:"profileType"`
		}{
			ID:          post.OrgID,
			Name:        post.OrgName,
			ProfileType: post.ProfileType,
		},
	}
}
