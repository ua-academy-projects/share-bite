package code

type Code string

var (
	NotFound     Code = "NOT_FOUND"
	BadRequest   Code = "BAD_REQUEST"
	Forbidden    Code = "FORBIDDEN"
	Unauthorized Code = "UNAUTHORIZED"
)
