package key

import (
	"fmt"
	"path"
	"strconv"
)

const (
	prefixCustomers  = "customers"
	prefixBusinesses = "businesses"
	prefixPosts      = "posts"
	prefixAvatars    = "avatars"
	prefixThumbnails = "thumbnails"
)

// CustomerPostImageKey generates an object key for a customer's post image.
// Result: customers/{customerID}/posts/{uploadSessionID}/{fileID}.{ext}
func CustomerPostImageKey(customerID, uploadSessionID, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixCustomers, customerID, prefixPosts, uploadSessionID, fileName)
}

// BusinessPostImageKey generates an object key for a business post image.
// Result: businesses/{unitID}/posts/{uploadSessionID}/{fileID}.{ext}
func BusinessPostImageKey(unitID int, uploadSessionID, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixBusinesses, strconv.Itoa(unitID), prefixPosts, uploadSessionID, fileName)
}

// CustomerAvatarKey generates an object key for a customer's avatar.
// Result: customers/{customerID}/avatars/{fileID}.{ext}
func CustomerAvatarKey(customerID, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixCustomers, customerID, prefixAvatars, fileName)
}

// getFileName constructs a safe filename by combining the file ID and extension.
func getFileName(fileID, ext string) string {
	return fmt.Sprintf("%s.%s", fileID, ext)
}

// PostThumbnailKey generates an object key for processed post thumbnails.
// Result: posts/thumbnails/{uploadSessionID}/{fileID}.{ext}
func PostThumbnailKey(uploadSessionID, fileID, ext string) string {
	fileName := getFileName(fileID, ext)

	return path.Join(
		prefixPosts,
		prefixThumbnails,
		uploadSessionID,
		fileName,
	)
}
