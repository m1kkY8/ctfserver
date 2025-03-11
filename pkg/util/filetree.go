package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateTree(root string, prefix string) string {
	var output strings.Builder

	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Sprintf("%sError reading directory: %v\n", prefix, err)
	}

	// Iterate over entries
	for i, entry := range entries {
		// Choose connector based on position.
		connector := "├── "
		newPrefix := prefix + "│   "
		if i == len(entries)-1 {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		output.WriteString(prefix + connector + entry.Name() + "\n")

		// If the entry is a directory, recursively add its contents.
		if entry.IsDir() {
			subDir := filepath.Join(root, entry.Name())
			output.WriteString(GenerateTree(subDir, newPrefix))
		}
	}
	return output.String()
}
