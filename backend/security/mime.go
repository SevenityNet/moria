package security

import (
	"moria/config"
)

var allowedMimeTypes []string
var allMimeTypeAllowed bool

func init() {
	allowedMimeTypes = config.GetSecurityAllowedMimeTypes()

	if len(allowedMimeTypes) == 0 || allowedMimeTypes[0] == "*" || allowedMimeTypes[0] == "images/*" {
		allMimeTypeAllowed = true
	}
}

func IsMimeTypeAllowed(mimeType string) bool {
	if allMimeTypeAllowed {
		return true
	}

	for _, allowedMimeType := range allowedMimeTypes {
		if allowedMimeType == mimeType {
			return true
		}
	}

	return false
}
