package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	finalDir   = "/var/encrypted_files"
	uploadDir  = "./uploads"
	publicDir  = "./public"
	serverPort = ":8080"
)

type UploadResponse struct {
	DownloadURL string `json:"download_url"`
	Hash        string `json:"hash"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request size
	r.Body = http.MaxBytesReader(w, r.Body, 200<<30) // 2GB

	err := r.ParseMultipartForm(100 << 20) // 100MB memory buffer
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		log.Println("Form parse error:", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not provided", http.StatusBadRequest)
		log.Println("File read error:", err)
		return
	}
	defer file.Close()

	// Generate unique filename
	ext := filepath.Ext(handler.Filename)
	b := make([]byte, 12) // 12 bytes = 16 chars base64url
	_, err = rand.Read(b)
	if err != nil {
		http.Error(w, "Random gen failed", http.StatusInternalServerError)
		return
	}
	shortName := base64.RawURLEncoding.EncodeToString(b)
	uniqueName := shortName + ext
	tmpPath := filepath.Join(uploadDir, uniqueName)
	publicPath := filepath.Join(publicDir, uniqueName)

	// Create temp file and hash
	outFile, err := os.Create(tmpPath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		log.Println("File create error:", err)
		return
	}
	defer outFile.Close()

	hash := sha256.New()
	multi := io.MultiWriter(outFile, hash)

	if _, err := io.Copy(multi, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		log.Println("Copy error:", err)
		return
	}

	// Move file to public dir
	err = os.Rename(tmpPath, publicPath)
	if err != nil {
		http.Error(w, "File move error", http.StatusInternalServerError)
		log.Println("Rename error:", err)
		return
	}

	// Build response
	fileHash := hex.EncodeToString(hash.Sum(nil))
	downloadURL := fmt.Sprintf("http://%s/download/%s", r.Host, uniqueName)

	log.Printf("Uploaded %s, hash=%s", uniqueName, fileHash)

	resp := UploadResponse{
		DownloadURL: downloadURL,
		Hash:        fileHash,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	// Ensure folders exist
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(publicDir, 0755)

	http.HandleFunc("/upload", uploadHandler)

	// Serve static files from /public via /download/
	fs := http.StripPrefix("/download/", http.FileServer(http.Dir(publicDir)))
	http.Handle("/download/", fs)

	log.Printf("Server listening on http://localhost%s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
