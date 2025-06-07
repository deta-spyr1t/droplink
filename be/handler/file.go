package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"droplink-be/model"
	"droplink-be/storage"
	"droplink-be/utils"
)

var backend storage.Storage

func Init(s storage.Storage) {
	backend = s
}

func Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	id := utils.GenerateID(12)
	storageKey := id + "_" + fileHeader.Filename

	_, err = backend.Save(file, storageKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	meta := &model.FileMeta{
		ID:           id,
		Filename:     fileHeader.Filename,
		Path:         storageKey,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
		MaxDownloads: 1,
		Downloads:    0,
	}
	model.Files[id] = meta

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"filename": fileHeader.Filename,
		"link":     fmt.Sprintf("/download/%s", id),
	})
}

func Download(c *gin.Context) {
	id := c.Param("id")
	meta, exists := model.Files[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if time.Now().After(meta.ExpiresAt) || meta.Downloads >= meta.MaxDownloads {
		_ = backend.Delete(meta.Path)
		delete(model.Files, id)
		c.JSON(http.StatusGone, gin.H{"error": "File expired or download limit reached"})
		return
	}

	data, filename, err := backend.Load(meta.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load file"})
		return
	}

	meta.Downloads++

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/octet-stream", data)
}
