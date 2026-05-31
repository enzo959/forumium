package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxFileSize = 20 << 20

var allowedTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
}

func UploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFileSize)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or oversized file (max 20 MB)"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized file type (PNG, JPEG, GIF only)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := time.Now().Format("20060102") + "_" + uuid.New().String() + ext
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating the uploads folder"})
		return
	}

	destPath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(header, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}

	url := "uploads/" + filename
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}