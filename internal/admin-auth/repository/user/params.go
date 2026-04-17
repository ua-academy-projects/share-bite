package user
type CreateUser struct {
	Email        string
	PasswordHash string
}

type CreatedUser struct {
	ID    string
	Email string
}