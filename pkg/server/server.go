package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/m1kkY8/ctfserver/pkg/config"
	"github.com/m1kkY8/ctfserver/pkg/handlers"
	"github.com/m1kkY8/ctfserver/pkg/logger"
	"github.com/m1kkY8/ctfserver/pkg/service"
)

const version = "1.0.0"

// Server represents the HTTP server
type Server struct {
	config      *config.Config
	httpServer  *http.Server
	fileService *service.FileService
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	fileService := service.NewFileService(cfg.RootDir, cfg.UploadDir, cfg.MaxUploadSize)

	return &Server{
		config:      cfg,
		fileService: fileService,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Initialize logger
	logger.InitLogger(s.config.LogLevel)

	// Create router
	router := s.setupRoutes()

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Logger.WithField("addr", addr).Info("Starting CTF file server")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	<-stop
	logger.Logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Logger.WithError(err).Error("Server forced to shutdown")
		return err
	}

	logger.Logger.Info("Server exited")
	return nil
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add middleware
	router.Use(logger.RecoveryMiddleware)
	router.Use(logger.LoggingMiddleware)

	// API routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	healthHandler := handlers.NewHealthHandler(version)
	apiRouter.Handle("/health", healthHandler).Methods("GET")

	// File tree endpoint
	fileTreeHandler := handlers.NewFileTreeHandler(s.fileService)
	apiRouter.Handle("/filetree", fileTreeHandler).Methods("GET")

	// Pretty file tree endpoint (human-readable, defaults to plain text)
	prettyFileTreeHandler := handlers.NewPrettyFileTreeHandler(s.fileService)
	apiRouter.Handle("/filetree/pretty", prettyFileTreeHandler).Methods("GET")

	// Shorter aliases for convenience
	apiRouter.Handle("/tree", prettyFileTreeHandler).Methods("GET") // Short alias for pretty tree
	apiRouter.Handle("/ls", prettyFileTreeHandler).Methods("GET")   // Unix-style alias

	// Upload endpoint
	uploadHandler := handlers.NewUploadHandler(s.fileService)
	apiRouter.Handle("/upload", uploadHandler).Methods("POST")

	// Uploads list endpoint
	uploadsListHandler := handlers.NewUploadsListHandler(s.fileService)
	apiRouter.Handle("/uploads", uploadsListHandler).Methods("GET")

	// Short alias for uploads list
	apiRouter.Handle("/ul", uploadsListHandler).Methods("GET")   // Short alias for uploads list
	apiRouter.Handle("/loot", uploadsListHandler).Methods("GET") // Short alias for uploads list

	// Static file server for downloads
	fs := http.StripPrefix("/files/", http.FileServer(http.Dir(s.config.RootDir)))
	router.PathPrefix("/files/").Handler(fs)

	return router
}
