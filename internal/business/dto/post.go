package dto

type PostResponse struct {
	ID       int64  `json:"id"`
	Content  string `json:"content"`
	ImageURL string `json:"imageUrl"`
}
