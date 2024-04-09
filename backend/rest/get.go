package rest

import (
	"log"
	"moria/processing"
	"moria/source"
	"strings"

	"github.com/gin-gonic/gin"
)

func getImage(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(400, gin.H{"error": "Category is missing"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID is missing"})
		return
	}

	src, err := source.GetCurrent().Get(category + "_" + id)
	if err != nil {
		log.Println(err)

		c.JSON(500, gin.H{"error": "Failed to get image"})
		return
	}

	ext := getExtension(id)
	if ext == "" {
		c.JSON(400, gin.H{"error": "ID is missing extension"})
		return
	}

	mimeType := getMimeType(ext)
	if mimeType == "" {
		c.JSON(400, gin.H{"error": "Unsupported file type"})
		return
	}

	if mimeType != "image/gif" {
		processed, err := processing.PostProcess(src, c)
		if err != nil {
			log.Println(err)

			c.JSON(500, gin.H{"error": "Failed to process image"})
			return
		}

		src = processed
	}

	c.Data(200, mimeType, src)
}

func getExtension(id string) string {
	split := strings.Split(id, ".")
	if len(split) < 2 {
		return ""
	}

	return split[len(split)-1]
}

func getMimeType(ext string) string {
	switch ext {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "webp":
		return "image/webp"
	case "gif":
		return "image/gif"
	default:
		return "image/jpeg"
	}
}
