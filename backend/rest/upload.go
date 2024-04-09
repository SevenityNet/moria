package rest

import (
	"io"
	"log"
	"moria/config"
	"moria/processing"
	"moria/security"
	"moria/source"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func postImage(c *gin.Context) {
	if !authenticate(c) {
		return
	}

	category := c.Param("category")
	if category == "" {
		c.JSON(400, gin.H{"error": "Category is missing"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File is missing"})
		return
	}

	mimeType := file.Header.Get("Content-Type")
	if !security.IsMimeTypeAllowed(mimeType) {
		c.JSON(400, gin.H{"error": "Mime type is not allowed"})
		return
	}

	src, err := file.Open()
	if err != nil {
		log.Println(err)

		c.JSON(500, gin.H{"error": "Failed to open file"})
		return
	}

	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		log.Println(err)

		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	ext := getExtensionWithDot(mimeType)
	if ext == "" {
		c.JSON(400, gin.H{"error": "Unsupported file type"})
		return
	}

	if ext != ".gif" {
		processedData, err := processing.Compress(data)
		if err != nil {
			log.Println(err)

			c.JSON(500, gin.H{"error": "Failed to compress file"})
			return
		}

		data = processedData
	}

	imageID := strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "") + ext
	filename := category + "_" + imageID

	err = source.GetCurrent().Upload(data, filename)
	if err != nil {
		log.Println(err)

		c.JSON(500, gin.H{"error": "Failed to upload file"})
		return
	}

	c.String(201, imageID)
}

func authenticate(c *gin.Context) bool {
	token := c.GetHeader(config.GetSecurityAPIAuthHeader())
	if token == "" {
		c.AbortWithStatus(401)
		return false
	}

	if token != config.GetSecurityAPIAuthToken() {
		c.AbortWithStatus(403)
		return false
	}

	return true
}

func getExtensionWithDot(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}
