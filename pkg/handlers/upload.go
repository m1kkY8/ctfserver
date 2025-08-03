package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/m1kkY8/ctfserver/pkg/logger"
	"github.com/m1kkY8/ctfserver/pkg/models"
	"github.com/m1kkY8/ctfserver/pkg/service"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	fileService *service.FileService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(fileService *service.FileService) *UploadHandler {
	return &UploadHandler{
		fileService: fileService,
	}
}

// ServeHTTP handles the file upload request
func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(h.fileService.MaxSize()); err != nil {
		logger.Logger.WithError(err).Error("Failed to parse multipart form")
		h.writeErrorResponse(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to get file from form")
		h.writeErrorResponse(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload the file
	result, err := h.fileService.UploadFile(header, file)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to upload file")
		h.writeErrorResponse(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Check if upload was successful
	if !result.Success {
		h.writeJSONResponse(w, result, http.StatusBadRequest)
		return
	}

	// Log successful upload
	logger.Logger.WithFields(map[string]interface{}{
		"filename": result.Filename,
		"size":     result.Size,
		"path":     result.Path,
	}).Info("File uploaded successfully")

	h.writeJSONResponse(w, result, http.StatusCreated)
}

func (h *UploadHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Logger.WithError(err).Error("Failed to encode JSON response")
	}
}

func (h *UploadHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := &models.ErrorResponse{
		Success: false,
		Error:   message,
	}
	h.writeJSONResponse(w, response, statusCode)
}
