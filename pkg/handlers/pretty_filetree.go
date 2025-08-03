package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/m1kkY8/ctfserver/pkg/logger"
	"github.com/m1kkY8/ctfserver/pkg/models"
	"github.com/m1kkY8/ctfserver/pkg/service"
)

// PrettyFileTreeHandler handles requests for human-readable file tree information
type PrettyFileTreeHandler struct {
	fileService *service.FileService
}

// NewPrettyFileTreeHandler creates a new pretty file tree handler
func NewPrettyFileTreeHandler(fileService *service.FileService) *PrettyFileTreeHandler {
	return &PrettyFileTreeHandler{
		fileService: fileService,
	}
}

// ServeHTTP handles the pretty file tree request
func (h *PrettyFileTreeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if client wants JSON response (default is plain text)
	acceptHeader := r.Header.Get("Accept")
	wantsJSON := acceptHeader == "application/json" || r.URL.Query().Get("format") == "json"

	result, err := h.fileService.GetPrettyFileTree()
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to generate pretty file tree")
		h.writeErrorResponse(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	if !result.Success {
		h.writeErrorResponse(w, result.Error, http.StatusInternalServerError)
		return
	}

	// Return JSON if specifically requested
	if wantsJSON {
		h.writeJSONResponse(w, result, http.StatusOK)
		return
	}

	// Return plain text by default
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result.TreeString))
}

func (h *PrettyFileTreeHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Logger.WithError(err).Error("Failed to encode JSON response")
	}
}

func (h *PrettyFileTreeHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := &models.ErrorResponse{
		Success: false,
		Error:   message,
	}
	h.writeJSONResponse(w, response, statusCode)
}
