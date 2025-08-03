package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/m1kkY8/ctfserver/pkg/logger"
	"github.com/m1kkY8/ctfserver/pkg/models"
	"github.com/m1kkY8/ctfserver/pkg/service"
)

// UploadsListHandler handles requests for listing uploaded files
type UploadsListHandler struct {
	fileService *service.FileService
}

// NewUploadsListHandler creates a new uploads list handler
func NewUploadsListHandler(fileService *service.FileService) *UploadsListHandler {
	return &UploadsListHandler{
		fileService: fileService,
	}
}

// ServeHTTP handles the uploads list request
func (h *UploadsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if client wants JSON response (default is plain text)
	acceptHeader := r.Header.Get("Accept")
	wantsJSON := acceptHeader == "application/json" || r.URL.Query().Get("format") == "json"

	prettyText, result, err := h.fileService.GetPrettyUploadsList()
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to list uploads")
		h.writeErrorResponse(w, "Failed to list uploads", http.StatusInternalServerError)
		return
	}

	if !result.Success {
		h.writeErrorResponse(w, result.Error, http.StatusInternalServerError)
		return
	}

	// Return JSON if specifically requested
	if wantsJSON {
		// Sort files by modification time (newest first) for JSON response
		sort.Slice(result.Files, func(i, j int) bool {
			return result.Files[i].ModTime.After(result.Files[j].ModTime)
		})
		h.writeJSONResponse(w, result, http.StatusOK)
		return
	}

	// Return plain text by default
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(prettyText))
}

func (h *UploadsListHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Logger.WithError(err).Error("Failed to encode JSON response")
	}
}

func (h *UploadsListHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := &models.ErrorResponse{
		Success: false,
		Error:   message,
	}
	h.writeJSONResponse(w, response, statusCode)
}
