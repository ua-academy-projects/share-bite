package dto

import "time"

type CommentResponse struct {
	ID        int64     `json:"id" example:"1"`
	PostID    int64     `json:"postId" example:"42"`
	AuthorID  string    `json:"authorId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content   string    `json:"content" example:"Great post!"`
	CreatedAt time.Time `json:"createdAt" example:"2024-04-06T10:00:00Z"`
}

type AuthorInfo struct {
	ID        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string  `json:"username" example:"john_doe"`
	FirstName string  `json:"firstName" example:"John"`
	LastName  string  `json:"lastName" example:"Doe"`
	AvatarURL *string `json:"avatarUrl,omitempty" example:"https://cdn.example.com/avatar.png"`
}

type CommentWithAuthorResponse struct {
	ID        int64      `json:"id" example:"1"`
	PostID    int64      `json:"postId" example:"42"`
	Content   string     `json:"content" example:"Great post!"`
	CreatedAt time.Time  `json:"createdAt" example:"2024-04-06T10:00:00Z"`
	Author    AuthorInfo `json:"author"`
}

type CreateCommentResponse struct {
	Comment CommentResponse `json:"comment"`
}

type UpdateCommentResponse struct {
	Comment CommentResponse `json:"comment"`
}

type GetCommentsResponse struct {
	Comments []CommentWithAuthorResponse `json:"comments"`
}
