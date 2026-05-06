package mediatype

import "strings"

// ExtFromContentType returns the file extension for the given MIME type.
// The second return value is false if the content type is not supported.
func ExtFromContentType(contentType string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(contentType))

	switch normalized {
	case "image/jpeg", "image/jpg":
		return "jpg", true
	case "image/png":
		return "png", true
	case "image/webp":
		return "webp", true
	case "image/heic", "image/heif":
		return "heic", true
	case "image/gif":
		return "gif", true
	default:
		return "", false
	}
}
