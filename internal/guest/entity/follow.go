package entity

import "time"

type CustomerFollow struct {
	ID string

	FollowerCustomerID string
	FollowedCustomerID string

	CreatedAt time.Time
}
