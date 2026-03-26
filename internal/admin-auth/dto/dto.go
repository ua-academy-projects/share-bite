package dto

import "github.com/ua-academy-projects/share-bite/internal/admin-auth/entity"

type CreateWithRoleParams struct {
	Email        string
	PasswordHash string
	RoleID       int
}

type CreatedUser struct {
	ID    string
	Email string
}

type UserWithRole struct {
	entity.User
	RoleSlug string
}
