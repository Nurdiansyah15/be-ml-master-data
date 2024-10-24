package utils

import (
	"path/filepath"
	"strings"
)

// extractPublicIDFromURL mengambil public ID dari URL Cloudinary
func ExtractPublicID(imageURL string) string {
	parts := strings.Split(imageURL, "/")
	filename := parts[len(parts)-1]
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
