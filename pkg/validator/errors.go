package validator

type ValidationError struct {
	Errors []ValidationErrorItem
}

type ValidationErrorItem struct {
	Field   string
	Message string
}
