package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/m1kkY8/ctfserver/pkg/server"
	"github.com/m1kkY8/ctfserver/pkg/util"
)

func fileTreeHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Start with the base directory name.
		baseName := filepath.Base(path)
		output := baseName + "\n"
		output += util.GenerateTree(path, "")

		// Return as plain text.
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, output)
	}
}

// GET request to download file specified file

func main() {
	// lhost := flag.String("lhost", "0.0.0.0", "IP address to listen on")
	lport := flag.String("lport", "8080", "Port to listen on")
	root := flag.String("root", ".", "Root directory to serve")
	flag.Parse()

	// parse port for server
	port := ":" + *lport

	fs := http.FileServer(http.Dir(*root))
	http.Handle("/", fs)

	http.HandleFunc("/filetree", fileTreeHandler(*root))

	// fmt.Println("Server started on :8080")
	// if err := http.ListenAndServe(port, nil); err != nil {
	// 	fmt.Printf("Server failed: %v\n", err)
	// }

	server.StartServer(port)
}
