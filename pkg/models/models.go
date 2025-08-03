package models

import "time"

// FileInfo represents a file or directory in the file tree
type FileInfo struct {
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"is_dir"`
	Size     int64      `json:"size,omitempty"`
	ModTime  time.Time  `json:"mod_time"`
	Children []FileInfo `json:"children,omitempty"`
}

// FileTreeResponse represents the response for file tree API
type FileTreeResponse struct {
	Success bool     `json:"success"`
	Root    FileInfo `json:"root,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// PrettyFileTreeResponse represents the response for human-readable file tree API
type PrettyFileTreeResponse struct {
	Success    bool     `json:"success"`
	Root       FileInfo `json:"root,omitempty"`
	TreeString string   `json:"tree_string,omitempty"`
	Error      string   `json:"error,omitempty"`
}

// UploadResponse represents the response for upload API
type UploadResponse struct {
	Success  bool   `json:"success"`
	Filename string `json:"filename,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Path     string `json:"path,omitempty"`
	Error    string `json:"error,omitempty"`
}

// UploadedFileInfo represents information about an uploaded file
type UploadedFileInfo struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	SizeHuman string    `json:"size_human"`
}

// UploadsListResponse represents the response for uploads list API
type UploadsListResponse struct {
	Success bool               `json:"success"`
	Files   []UploadedFileInfo `json:"files,omitempty"`
	Count   int                `json:"count"`
	Error   string             `json:"error,omitempty"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
