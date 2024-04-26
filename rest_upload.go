package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	extension, result, err := convertFile(data, fileType, currExt)
	if err != nil {
		panic(err)
	}

	uniqueFileID := strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", "") + "." + extension

	if result.OutputData != nil {
		if err := writeFileToDisk(subfolder, uniqueFileID, result.OutputData); err != nil {
			panic(err)
		}
	} else {
		if err := moveTmpToUploads(result.TmpFile, subfolder, uniqueFileID); err != nil {
			panic(err)
		}
	}

	c.JSON(201, gin.H{
		"file": uniqueFileID,
	})
}
