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
	prefixBanners    = "banners"
)

func CustomerPostImageKey(customerID, uploadSessionID, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixCustomers, customerID, prefixPosts, uploadSessionID, fileName)
}

func BusinessPostImageKey(unitID int, uploadSessionID, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixBusinesses, strconv.Itoa(unitID), prefixPosts, uploadSessionID, fileName)
}

func BusinessAvatarKey(unitID int, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixBusinesses, strconv.Itoa(unitID), prefixAvatars, fileName)
}

func BusinessBannerKey(unitID int, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixBusinesses, strconv.Itoa(unitID), prefixBanners, fileName)
}

func CustomerAvatarKey(customerID, fileID, ext string) string {
	fileName := getFileName(fileID, ext)
	return path.Join(prefixCustomers, customerID, prefixAvatars, fileName)
}

// getFileName constructs a safe filename by combining the file ID and extension.
func getFileName(fileID, ext string) string {
	return fmt.Sprintf("%s.%s", fileID, ext)
}
