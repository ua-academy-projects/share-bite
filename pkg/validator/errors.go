package validator

type ValidationError struct {
	Errors []ValidationErrorItem `json:"errors"`
}

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
