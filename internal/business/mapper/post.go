package mapper

import (
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
)

func ToPostResponse(post *entity.PostWithPhotos) dto.PostResponse {
	response := dto.PostResponse{
		ID:        post.ID,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		Images:    post.Images,
	}
	response.Org.ID = post.OrgID
	response.Org.Name = post.OrgName
	response.Org.ProfileType = post.ProfileType
	response.Org.Status = string(post.OrgStatus)

	return response
}
