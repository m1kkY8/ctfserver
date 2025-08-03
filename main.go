package main

import (
	"log"

	"github.com/m1kkY8/ctfserver/pkg/config"
	"github.com/m1kkY8/ctfserver/pkg/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create and start server
	srv := server.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
