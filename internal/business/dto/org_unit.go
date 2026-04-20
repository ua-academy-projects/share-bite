package dto

type CreateOrgRequest struct {
	ProfileType string   `json:"type" binding:"required"`
	ParentID    *int     `json:"parent_id,omitempty"`
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Latitude    *float32 `json:"latitude"`
	Longitude   *float32 `json:"longitude"`
}

type UpdateOrgRequest struct {
	Name        *string  `json:"name"`
	Avatar      *string  `json:"avatar"`
	Banner      *string  `json:"banner"`
	Description *string  `json:"description"`
	Latitude    *float32 `json:"latitude"`
	Longitude   *float32 `json:"longitude"`
}

type UpdateOrgResponse struct {
	Id          int      `json:"id" example:"42"`
	Name        string   `json:"name" example:"ShareBite Downtown"`
	Avatar      *string  `json:"avatar" example:"https://cdn.example.com/avatar.png"`
	Banner      *string  `json:"banner" example:"https://cdn.example.com/banner.png"`
	Description *string  `json:"description" example:"A cozy place in the city center."`
	Latitude    *float32 `json:"latitude" example:"50.4501"`
	Longitude   *float32 `json:"longitude" example:"30.5234"`
}

type UpdatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

