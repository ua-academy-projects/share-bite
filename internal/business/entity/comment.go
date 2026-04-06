package entity

import "time"

type Comment struct {
	ID        int64
	PostID    int64
	AuthorID  string
	Content   string
	CreatedAt time.Time
}

type CommentWithAuthor struct {
	ID        int64
	PostID    int64
	Content   string
	CreatedAt time.Time

	AuthorID        string
	AuthorUsername  string
	AuthorFirstName string
	AuthorLastName  string
	AuthorAvatarURL *string
}
