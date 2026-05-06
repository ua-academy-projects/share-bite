package response

type ErrorResponse struct {
	Message string        `json:"message" example:"error description or reason"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty" example:"field_name"`
	Message string `json:"message" example:"specific validation rule failed"`
}

type AuthErrorResponse struct {
	Error string `json:"error" example:"invalid or expired token"`
}
