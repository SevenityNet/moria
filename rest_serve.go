package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	file, err := getCachedFileIfExists(filePath)
	if err == ErrFileNotFound {
		log.Println("file not cached, saving and serving")
		file, err = os.ReadFile(filePath)
		if err != nil {
			panic(err)
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

	fileType, ok := getFileTypeByFilePath(filePath)
	if !ok {
		panic("invalid file type")
	}

	c.Data(http.StatusOK, FILETYPE_TO_CONTENTTYPE[fileType], file)
}
