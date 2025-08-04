.PHONY: build run test clean help deps

# Build the application
build:
	go build -o ctfserver

# Run the application with default settings
run:
	go run main.go

# Run with custom settings for development
dev:
	go run main.go -port 8080 -root ./testfiles -upload-dir ./uploads -log-level debug

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f ctfserver
	rm -rf uploads/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Create test directories
setup-dev:
	mkdir -p testfiles uploads
	echo "This is a test file" > testfiles/readme.txt
	echo "Another test file" > testfiles/example.txt

# Build for different platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o ctfserver-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o ctfserver-windows-amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o ctfserver-darwin-amd64

# Docker commands
docker-setup:
	./setup.sh

docker-build:
	UID=$$(id -u) GID=$$(id -g) docker-compose build

docker-up:
	UID=$$(id -u) GID=$$(id -g) docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f ctfserver

docker-test:
	./test_docker.sh

docker-clean:
	docker-compose down -v
	docker image prune -f

help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application with default settings"
	@echo "  dev        - Run with development settings"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  setup-dev  - Create test directories and files"
	@echo "  build-all  - Build for multiple platforms"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-setup  - Setup /opt directories with proper permissions"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start container"
	@echo "  docker-down   - Stop container"
	@echo "  docker-logs   - View container logs"
	@echo "  docker-test   - Test Docker setup"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  help          - Show this help message"
