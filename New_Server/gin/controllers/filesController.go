package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UploadHandler(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File not found",
		})
		return

	}

	// Case of dir to upload
	var saveDir string
	ext := filepath.Ext(file.Filename)
	if ext == ".txt" || ext == ".q" {
		saveDir = "uploads"
	} else {
		saveDir = "unknown"
	}

	// Make dir if not exist
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error creating folder",
		})
		fmt.Printf("Error creating folder: %s", err)
		return
	}

	// Get Filename without extension
	filename := strings.TrimSuffix(file.Filename, ext)

	// Make name of file with sufix "_1", "_2" etc, if file already exist
	newFilename := filepath.Join(saveDir, file.Filename)
	i := 1
	for {
		if _, err := os.Stat(newFilename); os.IsNotExist(err) {
			break
		}
		newFilename = filepath.Join(saveDir, fmt.Sprintf("%s_%d%s", filename, i, ext))
		i++
	}

	// Display file information
	fmt.Printf("Uploaded file: %s\n", file.Filename)
	fmt.Printf("File size: %d\n", file.Size)
	fmt.Printf("File type: %s\n", file.Header)

	// Save the file to a specific location on the server

	if err := c.SaveUploadedFile(file, newFilename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to save file",
		})

		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully!", "path": newFilename,
	})
}
