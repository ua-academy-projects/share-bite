package dto

import "time"

type ToggleLikeResponse struct {
	Liked bool `json:"liked" example:"true"`
}

type LikeItem struct {
	ID         int64     `json:"id" example:"1"`
	PostID     int64     `json:"postId" example:"42"`
	CustomerID string    `json:"customerId" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt  time.Time `json:"createdAt" example:"2024-04-06T10:00:00Z"`
}

type GetLikesResponse struct {
	Likes []LikeItem `json:"likes"`
}
