package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m1kkY8/ctfserver/pkg/models"
)

// GenerateFileTree creates a file tree structure starting from the given root path
func GenerateFileTree(root string) (*models.FileInfo, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	fileInfo := &models.FileInfo{
		Name:    info.Name(),
		Path:    root,
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(root)
		if err != nil {
			return fileInfo, nil // Return partial info even if can't read directory
		}

		for _, entry := range entries {
			childPath := filepath.Join(root, entry.Name())
			child, err := GenerateFileTree(childPath)
			if err != nil {
				continue // Skip files we can't read
			}
			fileInfo.Children = append(fileInfo.Children, *child)
		}
	}

	return fileInfo, nil
}

// GeneratePrettyTree creates a human-readable tree string from FileInfo
func GeneratePrettyTree(root *models.FileInfo) string {
	var builder strings.Builder
	builder.WriteString(root.Name)
	if root.IsDir {
		builder.WriteString("/")
	}
	builder.WriteString("\n")

	if root.Children != nil {
		generatePrettyTreeRecursive(root.Children, "", &builder)
	}

	return builder.String()
}

func generatePrettyTreeRecursive(children []models.FileInfo, prefix string, builder *strings.Builder) {
	for i, child := range children {
		isLast := i == len(children)-1

		// Choose the appropriate connector
		connector := "├── "
		newPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		// Write the current item
		builder.WriteString(prefix + connector + child.Name)
		if child.IsDir {
			builder.WriteString("/")
		} else {
			// Add file size for regular files
			builder.WriteString(fmt.Sprintf(" (%s)", formatFileSize(child.Size)))
		}
		builder.WriteString("\n")

		// Recursively handle children
		if child.Children != nil {
			generatePrettyTreeRecursive(child.Children, newPrefix, builder)
		}
	}
}

// formatFileSize converts bytes to human-readable format
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// IsValidFilename checks if a filename is safe (prevents path traversal)
func IsValidFilename(filename string) bool {
	if filename == "" || filename == "." || filename == ".." {
		return false
	}

	// Check for path separators that could indicate traversal attempts
	clean := filepath.Clean(filename)
	return clean == filename && !filepath.IsAbs(filename)
}
