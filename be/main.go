package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"droplink-be/handler"
	"droplink-be/storage"
	"droplink-be/storage/local"
	"droplink-be/storage/s3"
)

var storageBackend storage.Storage

func initStorage() {
	driver := os.Getenv("STORAGE_DRIVER")

	switch driver {
	case "s3":
		bucket := os.Getenv("S3_BUCKET")
		region := os.Getenv("AWS_REGION")

		s3Storage, err := s3.NewS3Storage(region, bucket)
		if err != nil {
			log.Fatalf("Failed to initialize S3 storage: %v", err)
		}
		storageBackend = s3Storage
	case "local":
		path := os.Getenv("LOCAL_STORAGE_PATH")
		if path == "" {
			path = "uploads"
		}
		storageBackend = local.NewLocalStorage(path)
	default:
		log.Fatalf("Unsupported STORAGE_DRIVER: %s", driver)
	}
}

func main() {
	initStorage()

	r := gin.Default()

	// Inject storage backend into handlers
	handler.Init(storageBackend)

	r.POST("/upload", handler.Upload)
	r.GET("/download/:id", handler.Download)

	r.Run(":8080")
}
