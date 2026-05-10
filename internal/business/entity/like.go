package entity

import "time"

type Like struct {
	ID        int64
	PostID    int64
	AuthorID  string
	CreatedAt time.Time
}

type LikeWithAuthor struct {
	ID        int64
	PostID    int64
	CreatedAt time.Time

	AuthorID        string
	AuthorUsername  string
	AuthorFirstName string
	AuthorLastName  string
	AuthorAvatarURL *string
}
