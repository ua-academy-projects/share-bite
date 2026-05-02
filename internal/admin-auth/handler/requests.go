package handler

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Slug     string `json:"slug" binding:"required,oneof=user business"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type OAuthCallbackRequest struct {
	Code string `json:"code" binding:"required"`
	Slug string `json:"slug" binding:"required,oneof=user business"`
}

type OAuthLinkRequest struct {
	Code string `json:"code" binding:"required"`
}

type RecoverAccessRequest struct {
	Email string `json:"email" binding:"required,email,max=254"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UsersFilterQuery struct {
	Limit     *int   `form:"limit" binding:"omitempty,min=1"`
	Offset    *int   `form:"offset" binding:"omitempty,min=0"`
	Search    string `form:"search_email" binding:"omitempty,max=255" `
	Role      string `form:"role" binding:"omitempty,oneof=admin moderator user business"`
	Status    string `form:"status" binding:"omitempty,oneof=active muted suspended"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

type ChangeRoleRequest struct {
	RoleSlug string `json:"role_slug" binding:"required,oneof=admin moderator user business"`
}
