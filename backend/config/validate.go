package config

import (
	"fmt"
	"strings"
)

func Validate() error {
	cacheEnabled, err := getBoolEnv("CACHE_ENABLED")
	if err != nil {
		return err
	}
	fmt.Println("CACHE_ENABLED:", cacheEnabled)

	cacheExpiration, err := getIntEnv("CACHE_EXPIRATION")
	if err != nil {
		return err
	}
	fmt.Println("CACHE_EXPIRATION:", cacheExpiration)

	cacheHash, err := getStringEnv("CACHE_HASH", []string{"none", "sha256", "md5", "crc32"})
	if err != nil {
		return err
	}
	fmt.Println("CACHE_HASH:", cacheHash)

	cacheHashIncludesDimensions, err := getBoolEnv("CACHE_HASH_INCLUDES_DIMENSIONS")
	if err != nil {
		return err
	}
	fmt.Println("CACHE_HASH_INCLUDES_DIMENSIONS:", cacheHashIncludesDimensions)

	sourceType, err := getStringEnv("SOURCE_TYPE", []string{"local", "remote"})
	if err != nil {
		return err
	}
	fmt.Println("SOURCE_TYPE:", sourceType)

	if sourceType == "local" {
		_, err := getStringEnv("SOURCE_LOCAL_PATH", nil) // No enum validation required
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(sourceType, "remote") {
		_, err = getStringEnv("SOURCE_REMOTE_HOST", nil)
		if err != nil {
			return err
		}
		_, err = getIntEnv("SOURCE_REMOTE_PORT")
		if err != nil {
			return err
		}

		switch strings.TrimPrefix(sourceType, "remote_") {
		case "ftp", "sftp":
			_, err = getStringEnv("SOURCE_REMOTE_FTP_PATH", nil)
			if err != nil {
				return err
			}
			_, err = getStringEnv("SOURCE_REMOTE_FTP_USER", nil)
			if err != nil {
				return err
			}
			_, err = getStringEnv("SOURCE_REMOTE_FTP_PASS", nil)
			if err != nil {
				return err
			}
		case "ssh":
			_, err = getStringEnv("SOURCE_REMOTE_SSH_USER", nil)
			if err != nil {
				return err
			}
		case "upload":
			_, err = getStringEnv("SOURCE_REMOTE_UPLOAD_URL", nil)
			if err != nil {
				return err
			}
			uploadAuth, err := getStringEnv("SOURCE_REMOTE_UPLOAD_AUTH", []string{"none", "custom", "basic", "bearer"})
			if err != nil {
				return err
			}
			fmt.Println("SOURCE_REMOTE_UPLOAD_AUTH:", uploadAuth)
		}
	}

	processingEnabled, err := getBoolEnv("PROCESSING_ENABLED")
	if err != nil {
		return err
	}
	fmt.Println("PROCESSING_ENABLED:", processingEnabled)

	if processingEnabled {
		_, err = getBoolEnv("PROCESSING_COMPRESSION_ENABLED")
		if err != nil {
			return err
		}
		_, err = getBoolEnv("PROCESSING_COMPRESSION_LOSSLESS")
		if err != nil {
			return err
		}
		_, err = getBoolEnv("PROCESSING_RESIZE_ENABLED")
		if err != nil {
			return err
		}
		_, err = getBoolEnv("PROCESSING_CROP_ENABLED")
		if err != nil {
			return err
		}
		_, err = getBoolEnv("PROCESSING_TO_GRAYSCALE_ENABLED")
		if err != nil {
			return err
		}
	}

	_, err = getBoolEnv("API_ENABLED")
	if err != nil {
		return err
	}

	_, err = getBoolEnv("SECURITY_CORS_ENABLED")
	if err != nil {
		return err
	}

	_, err = getBoolEnv("FRONTEND_ENABLED")
	if err != nil {
		return err
	}

	return nil
}
