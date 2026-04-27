package entity

import "time"

type CustomerFollow struct {
	ID string

	FollowerCustomerID string
	FollowedCustomerID string

	CreatedAt time.Time
}

type Follower struct {
	Customer
	FollowCreatedAt time.Time
	FollowID        string

	IsFollowing  bool
	IsFollowedBy bool
	IsMutual     bool
}

type ListFollowersInput struct {
	TargetCustomerID    string
	RequesterCustomerID *string
	PageSize            int
	PageToken           string
}

type ListFollowersOutput struct {
	Followers     []Follower
	NextPageToken string
}

type ListFollowingInput struct {
	TargetCustomerID    string
	RequesterCustomerID *string
	PageSize            int
	PageToken           string
}

type ListFollowingOutput struct {
	Followers     []Follower
	NextPageToken string
}
