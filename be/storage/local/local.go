package local

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"droplink-be/storage"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(base string) storage.Storage {
	os.MkdirAll(base, os.ModePerm)
	return &LocalStorage{BasePath: base}
}

func (s *LocalStorage) Save(file multipart.File, filename string) (string, error) {
	dest := filepath.Join(s.BasePath, filename)
	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return filename, err
}

func (s *LocalStorage) Load(path string) ([]byte, string, error) {
	full := filepath.Join(s.BasePath, path)
	data, err := os.ReadFile(full) // ✅ modern replacement for ioutil.ReadFile
	return data, path, err
}

func (s *LocalStorage) Delete(path string) error {
	return os.Remove(filepath.Join(s.BasePath, path))
}
