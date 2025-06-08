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

	"github.com/rs/cors"
)

const (
	uploadDir     = "/tmp/uploads"
	publicDir     = "/tmp/public"
	serverPort    = ":8080"
	frontEndpoint = "http://localhost:3000"
	verbose       = false
)

type UploadResponse struct {
	DownloadURL string `json:"download_url"`
	Hash        string `json:"hash"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// Get IP address (considering proxy headers too)
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Get User-Agent
	userAgent := r.UserAgent()

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

	//Log more data
	if verbose {
		log.Printf("\nUploaded: %s\nIP: %s\nUser Agent: %s\nFile Hash: %s", uniqueName, ip, userAgent, fileHash)
	} else {
		log.Printf("\nUploaded: %s\nFile Hash: %s", uniqueName, fileHash)
	}
	resp := UploadResponse{
		DownloadURL: downloadURL,
		Hash:        fileHash,
	}
	w.Header().Set("Access-Control-Allow-Origin", frontEndpoint)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()

	// Ensure folders exist
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(publicDir, 0755)

	mux.HandleFunc("/upload", uploadHandler)

	// Serve static files from /public via /download/
	fs := http.StripPrefix("/download/", http.FileServer(http.Dir(publicDir)))
	mux.Handle("/download/", fs)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	log.Printf("Server listening on http://localhost%s", serverPort)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
