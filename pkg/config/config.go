package config

import (
	"flag"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Host          string
	Port          int
	RootDir       string
	UploadDir     string
	MaxUploadSize int64
	LogLevel      string
}

// LoadConfig loads configuration from environment variables and command line flags
func LoadConfig() *Config {
	cfg := &Config{
		Host:          getEnvOrDefault("CTF_HOST", "0.0.0.0"),
		Port:          getEnvOrDefaultInt("CTF_PORT", 8080),
		RootDir:       getEnvOrDefault("CTF_ROOT_DIR", "."),
		UploadDir:     getEnvOrDefault("CTF_UPLOAD_DIR", "./uploads"),
		MaxUploadSize: getEnvOrDefaultInt64("CTF_MAX_UPLOAD_SIZE", 200*1024*1024), // 200MB
		LogLevel:      getEnvOrDefault("CTF_LOG_LEVEL", "info"),
	}

	// Command line flags override environment variables
	host := flag.String("host", cfg.Host, "Host to bind to")
	port := flag.Int("port", cfg.Port, "Port to listen on")
	rootDir := flag.String("root", cfg.RootDir, "Root directory to serve")
	uploadDir := flag.String("upload-dir", cfg.UploadDir, "Directory for uploaded files")
	maxUpload := flag.Int64("max-upload", cfg.MaxUploadSize, "Maximum upload size in bytes")
	logLevel := flag.String("log-level", cfg.LogLevel, "Log level (debug, info, warn, error)")
	flag.Parse()

	cfg.Host = *host
	cfg.Port = *port
	cfg.RootDir = *rootDir
	cfg.UploadDir = *uploadDir
	cfg.MaxUploadSize = *maxUpload
	cfg.LogLevel = *logLevel

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvOrDefaultInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}
