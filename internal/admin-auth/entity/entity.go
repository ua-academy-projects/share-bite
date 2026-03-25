package entity

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type Role struct {
	ID   int
	Slug string
	Name string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
