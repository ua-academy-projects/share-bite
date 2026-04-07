package dto

type PostResponse struct {
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}
