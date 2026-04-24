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
}

type ListFollowersInput struct {
	TargetCustomerID    string
	RequesterCustomerID *string
	PageSize            int
	PageToken           string
}

type ListFollowersOutput struct {
	Customers     []Customer
	NextPageToken string
}

type ListFollowingInput struct {
	TargetCustomerID    string
	RequesterCustomerID *string
	PageSize            int
	PageToken           string
}

type ListFollowingOutput struct {
	Customers     []Customer
	NextPageToken string
}
