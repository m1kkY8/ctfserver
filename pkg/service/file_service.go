package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/m1kkY8/ctfserver/pkg/models"
	"github.com/m1kkY8/ctfserver/pkg/util"
)

// FileService handles file operations
type FileService struct {
	rootDir   string
	uploadDir string
	maxSize   int64
}

// NewFileService creates a new file service
func NewFileService(rootDir, uploadDir string, maxSize int64) *FileService {
	return &FileService{
		rootDir:   rootDir,
		uploadDir: uploadDir,
		maxSize:   maxSize,
	}
}

// MaxSize returns the maximum allowed file size
func (fs *FileService) MaxSize() int64 {
	return fs.maxSize
}

// GetFileTree returns the file tree for the root directory
func (fs *FileService) GetFileTree() (*models.FileInfo, error) {
	return util.GenerateFileTree(fs.rootDir)
}

// GetPrettyFileTree returns both structured and human-readable file tree
func (fs *FileService) GetPrettyFileTree() (*models.PrettyFileTreeResponse, error) {
	fileTree, err := util.GenerateFileTree(fs.rootDir)
	if err != nil {
		return &models.PrettyFileTreeResponse{
			Success: false,
			Error:   "Failed to read directory",
		}, nil
	}

	prettyTree := util.GeneratePrettyTree(fileTree)

	return &models.PrettyFileTreeResponse{
		Success:    true,
		Root:       *fileTree,
		TreeString: prettyTree,
	}, nil
}

// UploadFile saves an uploaded file to the upload directory
func (fs *FileService) UploadFile(fileHeader *multipart.FileHeader, file multipart.File) (*models.UploadResponse, error) {
	// Validate file size
	if fileHeader.Size > fs.maxSize {
		return &models.UploadResponse{
			Success: false,
			Error:   fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", fs.maxSize),
		}, nil
	}

	// Validate filename
	filename := filepath.Base(fileHeader.Filename)
	if !util.IsValidFilename(filename) {
		return &models.UploadResponse{
			Success: false,
			Error:   "Invalid filename",
		}, nil
	}

	// Ensure upload directory exists
	if err := util.EnsureDir(fs.uploadDir); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Create destination file
	dstPath := filepath.Join(fs.uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	written, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	return &models.UploadResponse{
		Success:  true,
		Filename: filename,
		Size:     written,
		Path:     dstPath,
	}, nil
}
