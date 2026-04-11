package code

type Code string

var (
	NotFound Code = "NOT_FOUND"

	AlreadyExists Code = "ALREADY_EXISTS"
	UpstreamError Code = "UPSTREAM_ERROR"

	InvalidJSON    Code = "INVALID_JSON"
	InvalidRequest Code = "INVALID_REQUEST"

	EmptyUpdate Code = "EMPTY_UPDATE"

	BadRequest Code = "BAD_REQUEST"
	Internal   Code = "INTERNAL"
)
