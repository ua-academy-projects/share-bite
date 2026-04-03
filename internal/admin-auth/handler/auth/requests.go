package auth

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
