package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {
	if !authenticateCodeOrToken(c) {
		return
	}

	subfolder := c.Param("folder")
	if subfolder == "" {
		c.JSON(400, gin.H{
			"error": "Subfolder is required",
		})

		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "File is required",
		})

		return
	}

	fileType, ok := getFileType(file)
	if !ok {
		c.JSON(400, gin.H{
			"error": "Invalid file type",
		})

		return
	}

	data := make([]byte, file.Size)
	fileData, err := file.Open()
	if err != nil {
		panic(err)
	}

	defer fileData.Close()

	_, err = fileData.Read(data)
	if err != nil {
		panic(err)
	}

	currExt := strings.Split(file.Filename, ".")[1]

	result, err := convertFile(data, fileType, subfolder, currExt)
	if err != nil {
		panic(err)
	}

	if result.OutputData != nil {
		if err := writeFileToDisk(subfolder, result.OutputFileID, result.OutputData); err != nil {
			panic(err)
		}
	} else if result.TmpFile != "" {
		if err := moveTmpToUploads(result.TmpFile, subfolder, result.OutputFileID); err != nil {
			panic(err)
		}
	}

	c.JSON(201, gin.H{
		"file": result.OutputFileID,
	})
}
