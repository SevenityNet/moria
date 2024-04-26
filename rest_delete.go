package main

import "github.com/gin-gonic/gin"

func deleteFile(c *gin.Context) {
	if !authenticateBackend(c) {
		return
	}

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

	if err := deleteFileFromDisk(subfolder, filename); err != nil {
		if err == ErrFileNotFound {
			c.JSON(404, gin.H{
				"error": "File not found",
			})

			return
		}

		panic(err)
	}

	c.Status(204)
}
