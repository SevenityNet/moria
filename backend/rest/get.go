package rest

import (
	"log"
	"moria/cache"
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

	cacheKey := cache.Key(category+"_"+id, c.Request.URL.Query())
	existing := cache.Get(cacheKey)
	if existing != nil {
		c.Data(200, mimeType, existing)
		return
	}

	src, err := source.GetCurrent().Get(category + "_" + id)
	if err != nil {
		log.Println(err)

		c.JSON(500, gin.H{"error": "Failed to get image"})
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

	cache.Set(cacheKey, src)

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
