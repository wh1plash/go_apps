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

	// Get Filename without extension
	filename := strings.TrimSuffix(file.Filename, ext)

	// Make name of file with sufix "_1", "_2" etc, if file already exist
	newFilename, err := getUniqueFilename(saveDir, filename, ext)
	if err != nil {
		fmt.Printf("Ошибка при получении уникального имени файла: " + err.Error())
		return
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

func getUniqueFilename(saveDir, filename, ext string) (string, error) {
	newFilename := filepath.Join(saveDir, filename+ext)
	i := 1

	// Проверяем существование файла и генерируем новое имя, если файл существует
	for {
		if _, err := os.Stat(newFilename); os.IsNotExist(err) {
			return newFilename, nil // Возвращаем уникальное имя файла
		}
		newFilename = filepath.Join(saveDir, fmt.Sprintf("%s_%d%s", filename, i, ext))
		i++
	}
}
