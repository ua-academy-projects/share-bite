package dto

import customerResponse "github.com/ua-academy-projects/share-bite/internal/guest/handler/customer"

type ListFollowersRequest struct {
	PageSize  int    `form:"page_size,default=20" binding:"gte=1,lte=100"`
	PageToken string `form:"page_token"`
}

type ListFollowingRequest struct {
	PageSize  int    `form:"page_size,default=20" binding:"gte=1,lte=100"`
	PageToken string `form:"page_token"`
}

type ListCustomersResponse struct {
	Customers     []customerResponse.CustomerResponse `json:"customers"`
	NextPageToken string                              `json:"next_page_token,omitempty"`
}
