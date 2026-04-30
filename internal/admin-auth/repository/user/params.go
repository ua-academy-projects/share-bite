package user

import "github.com/ua-academy-projects/share-bite/internal/admin-auth/models"

type CreateUser struct {
	Email        string
	PasswordHash string
}

type CreatedUser struct {
	ID    string
	Email string
}

type UpdateUserStatus struct {
	UserID  string
	Status  models.UserStatus
	SetByID string
}
