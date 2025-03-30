package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	uploadDir     = "./exfil"
	maxUploadSize = 200 * 1024 * 1024 // 200 MB
)

func UploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit upload size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		// Parse multipart form
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		// Get file from form data
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Invalid file in form data", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Sanitize filename to prevent path traversal
		filename := filepath.Base(header.Filename)
		if filename == "." || filename == string(filepath.Separator) {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		// Create destination file
		dstPath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy file content
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// Log the upload
		fmt.Printf("[%s] Uploaded file: %s (%d bytes)\n",
			time.Now().Format(time.RFC3339),
			filename,
			header.Size,
		)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "File uploaded successfully: %s\n", filename)
	}
}
