package main

import (
	"flag"
	"net/http"

	"github.com/m1kkY8/ctfserver/pkg/api"
	"github.com/m1kkY8/ctfserver/pkg/server"
)

func main() {
	// lhost := flag.String("lhost", "0.0.0.0", "IP address to listen on")
	lport := flag.String("lport", "8080", "Port to listen on")
	root := flag.String("root", ".", "Root directory to serve")
	flag.Parse()

	// parse port for server
	port := ":" + *lport

	fs := http.StripPrefix("/", http.FileServer(http.Dir(*root)))
	http.Handle("/", fs)

	http.HandleFunc("/filetree", api.FileTreeHandler(*root))
	http.HandleFunc("/upload", api.UploadHandler())

	server.StartServer(port)
}
