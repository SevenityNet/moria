package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	FILETYPE_TO_CONTENTTYPE = map[fileType]string{
		IMAGE: "image/webp",
		VIDEO: "video/webm",
		AUDIO: "audio/mp3",
	}
)

func serveFile(c *gin.Context) {
	subfolder := c.Param("folder")
	if subfolder == "" {
		c.JSON(400, gin.H{
			"error": "Subfolder is required",
		})

		return
	}

	filename := c.Param("file")
	if filename == "" {
		c.JSON(400, gin.H{
			"error": "Filename is required",
		})

		return
	}

	filePath := getFilePath(subfolder, filename)

	fileType, ok := getFileTypeByFilePath(filePath)
	if !ok {
		panic("invalid file type")
	}

	file, err := getCachedFileIfExists(filePath)
	if err == ErrFileNotFound {
		log.Println("file not cached, saving and serving")
		file, err = os.ReadFile(filePath)
		if errors.Is(err, os.ErrNotExist) && fileType == VIDEO {
			file, ext := findFileWithDifferentExtension(subfolder, filename)
			if file == nil {
				c.JSON(404, gin.H{"error": "File not found"})
				return
			}

			contentType := fmt.Sprintf("video/%s", ext)
			c.Data(http.StatusOK, contentType, file)
			return
		} else if err != nil {
			c.JSON(404, gin.H{"error": "File not found"})
			return
		}

		err = cacheFile(filePath, file)
		if err != nil {
			fmt.Sprintln(err)
		}

	} else if err != nil {
		panic(err)
	} else {
		log.Println("Serving from cache")
	}

	if fileType == IMAGE {
		file, err = postprocessImage(c, file)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Failed to process image"})
			return
		}
	}

	c.Data(http.StatusOK, FILETYPE_TO_CONTENTTYPE[fileType], file)
}

// Helper Function to search for file with different video extensions
func findFileWithDifferentExtension(subfolder, filename string) ([]byte, string) {
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename)) // Remove original extension

	for key := range ALLOWED_VIDEO_EXTENSIONS {
		potentialPath := getFilePath(subfolder, fmt.Sprintf("%s.%s", baseName, key))
		file, err := os.ReadFile(potentialPath)
		if err == nil {
			return file, key
		}
	}

	return nil, "" // No file found
}
