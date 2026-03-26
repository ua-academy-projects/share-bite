package entity

import (
	"time"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type Role struct {
	ID   int    `db:"id" json:"id"`
	Slug string `db:"slug" json:"slug"`
	Name string `db:"name" json:"name"`
}

type UserRole struct {
	UserID string `db:"user_id" json:"user_id"`
	RoleID int    `db:"role_id" json:"role_id"`
}
