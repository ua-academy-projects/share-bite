package dto

import "time"

type PostResponse struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	Org struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		ProfileType string `json:"profileType"`
	} `json:"org"`

	Images []string `json:"images"`
}
