package entity

import "time"

type Customer struct {
	ID     string
	UserID string

	UserName  string
	FirstName string
	LastName  string

	AvatarObjectKey *string
	Bio             *string

	IsFollowersPublic bool
	IsFollowingPublic bool

	CreatedAt time.Time
}

type CreateCustomer struct {
	UserID string

	UserName  string
	FirstName string
	LastName  string

	Bio *string
}

type UpdateCustomer struct {
	UserID string

	UserName  *string
	FirstName *string
	LastName  *string

	AvatarObjectKey *string
	Bio             *string

	IsFollowersPublic *bool
	IsFollowingPublic *bool
}
