package entity

import "time"

type Collaborator struct {
	CollectionID string
	CustomerID   string

	UserName        string
	AvatarObjectKey *string

	AddedAt time.Time
}

type RemoveCollaboratorInput struct {
	CollectionID string

	CustomerID       string
	TargetCustomerID string
}
