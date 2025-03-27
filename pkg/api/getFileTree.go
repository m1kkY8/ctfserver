package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/m1kkY8/ctfserver/pkg/util"
)

func FileTreeHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Start with the base directory name.
		baseName := filepath.Base(path)
		output := baseName + "\n"
		output += util.GenerateTree(path, "")

		// Return as plain text.
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, output)

		// 127.0.0.1 - - [28/Mar/2025 00:17:05] "GET / HTTP/1.1" 200 -
		// log like the line from above
	}
}
