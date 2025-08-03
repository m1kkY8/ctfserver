package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/m1kkY8/ctfserver/pkg/logger"
	"github.com/m1kkY8/ctfserver/pkg/models"
	"github.com/m1kkY8/ctfserver/pkg/service"
)

// FileTreeHandler handles requests for file tree information
type FileTreeHandler struct {
	fileService *service.FileService
}

// NewFileTreeHandler creates a new file tree handler
func NewFileTreeHandler(fileService *service.FileService) *FileTreeHandler {
	return &FileTreeHandler{
		fileService: fileService,
	}
}

// ServeHTTP handles the file tree request
func (h *FileTreeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileTree, err := h.fileService.GetFileTree()
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to generate file tree")
		h.writeErrorResponse(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	response := &models.FileTreeResponse{
		Success: true,
		Root:    *fileTree,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

func (h *FileTreeHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Logger.WithError(err).Error("Failed to encode JSON response")
	}
}

func (h *FileTreeHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := &models.ErrorResponse{
		Success: false,
		Error:   message,
	}
	h.writeJSONResponse(w, response, statusCode)
}
