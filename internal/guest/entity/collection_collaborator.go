package entity

import "time"

type Collaborator struct {
	CollectionID string
	CustomerID   string

	UserName        string
	AvatarObjectKey *string

	CreatedAt time.Time
}

type RemoveCollaboratorInput struct {
	CollectionID string

	CustomerID       string
	TargetCustomerID string
}
