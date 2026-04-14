package ghAuth

import "time"

type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	ID        int64
	GitHubID  int64
	Login     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
