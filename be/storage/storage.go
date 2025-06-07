package storage

import (
	"mime/multipart"
)

type Storage interface {
	Save(file multipart.File, filename string) (string, error)
	Load(path string) ([]byte, string, error)
	Delete(path string) error
}
