package version

// These variables are overwritten during compilation via -ldflags.
// The default values are used when running `go run` locally.
var (
	Version    = "development"
	CommitHash = "unknown"
	BuildTime  = "unknown"
)
