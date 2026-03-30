package code

type Code string

var (
	NotFound Code = "NOT_FOUND"

	AlreadyExists Code = "ALREADY_EXISTS"

	InvalidJSON    Code = "INVALID_JSON"
	InvalidRequest Code = "INVALID_REQUEST"

	EmptyUpdate Code = "EMPTY_UPDATE"
)
