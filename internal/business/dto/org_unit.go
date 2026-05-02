package dto

type CreateOrgRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	Avatar      *string  `json:"avatar"` 
    Banner      *string  `json:"banner"`
}

type UpdateOrgRequest struct {
	Name        *string  `json:"name"`
	Avatar      *string  `json:"avatar"`
	Banner      *string  `json:"banner"`
	Description *string  `json:"description"`
}

type UpdateOrgResponse struct {
	Id          int      `json:"id" example:"42"`
	Name        string   `json:"name" example:"ShareBite Downtown"`
	Avatar      *string  `json:"avatar" example:"https://cdn.example.com/avatar.png"`
	Banner      *string  `json:"banner" example:"https://cdn.example.com/banner.png"`
	Description *string  `json:"description" example:"A cozy place in the city center."`
}

type UpdatePostRequest struct {
	Content string `json:"content" binding:"required"`
}
