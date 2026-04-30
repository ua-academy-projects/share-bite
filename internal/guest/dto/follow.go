package dto

type ListFollowersRequest struct {
	PageSize  int    `form:"page_size,default=20" binding:"gte=1,lte=100"`
	PageToken string `form:"page_token"`
}

type ListFollowingRequest struct {
	PageSize  int    `form:"page_size,default=20" binding:"gte=1,lte=100"`
	PageToken string `form:"page_token"`
}

type ListCustomersResponse struct {
	Customers     []FollowerResponse `json:"customers"`
	NextPageToken string             `json:"next_page_token,omitempty"`
}

type FollowerResponse struct {
	ID string `json:"id"`

	UserName  string  `json:"userName"`
	AvatarURL *string `json:"avatarUrl"`

	IsFollowing  bool `json:"isFollowing"`
	IsFollowedBy bool `json:"isFollowedBy"`
	IsMutual     bool `json:"isMutual"`
	IsMe         bool `json:"isMe"`
}
