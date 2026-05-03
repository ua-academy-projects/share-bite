package mediatype

import "strings"

// ExtFromContentType returns the file extension for supported MIME types.
// It returns an empty string if the content type is not supported.
func ExtFromContentType(contentType string) string {
	normalized := strings.ToLower(strings.TrimSpace(contentType))

	switch normalized {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/webp":
		return "webp"
	case "image/heic", "image/heif":
		return "heic"
	case "image/gif":
		return "gif"
	default:
		return ""
	}
}
