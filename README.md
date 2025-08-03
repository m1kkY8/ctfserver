# CTF File Server

A lightweight, secure file server designed for CTF (Capture The Flag) competitions to facilitate easy file transfers between machines.

## Features

- **File Upload**: Secure file upload with size limits and validation
- **File Download**: Static file serving for downloads
- **Directory Listing**: JSON API for browsing file structures
- **Health Checks**: Built-in health check endpoint
- **Structured Logging**: JSON-formatted logs with request tracking
- **Graceful Shutdown**: Proper server shutdown handling
- **Configurable**: Environment variables and command-line flag support

## Quick Start

### Build and Run

```bash
# Build the server
go build -o ctfserver

# Run with default settings
./ctfserver

# Run with custom configuration
./ctfserver -host 0.0.0.0 -port 8080 -root ./files -upload-dir ./uploads
```

### Configuration

The server can be configured via environment variables or command-line flags:

#### Environment Variables

- `CTF_HOST`: Host to bind to (default: "0.0.0.0")
- `CTF_PORT`: Port to listen on (default: 8080)
- `CTF_ROOT_DIR`: Root directory for file downloads (default: ".")
- `CTF_UPLOAD_DIR`: Directory for uploaded files (default: "./uploads")
- `CTF_MAX_UPLOAD_SIZE`: Maximum upload size in bytes (default: 209715200 = 200MB)
- `CTF_LOG_LEVEL`: Log level - debug, info, warn, error (default: "info")

#### Command-Line Flags

- `-host`: Host to bind to
- `-port`: Port to listen on
- `-root`: Root directory for file downloads
- `-upload-dir`: Directory for uploaded files
- `-max-upload`: Maximum upload size in bytes
- `-log-level`: Log level

## API Endpoints

### Health Check

```bash
GET /api/v1/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-08-04T10:30:00Z",
  "version": "1.0.0"
}
```

### File Tree

Get the directory structure as JSON:

```bash
GET /api/v1/filetree
```

Response:
```json
{
  "success": true,
  "root": {
    "name": "files",
    "path": "./files",
    "is_dir": true,
    "mod_time": "2025-08-04T10:30:00Z",
    "children": [
      {
        "name": "example.txt",
        "path": "./files/example.txt",
        "is_dir": false,
        "size": 1024,
        "mod_time": "2025-08-04T10:30:00Z"
      }
    ]
  }
}
```

### Pretty File Tree (Human-Readable)

Get a human-readable directory tree (defaults to plain text):

```bash
GET /api/v1/filetree/pretty
GET /api/v1/tree          # Short alias
GET /api/v1/ls            # Unix-style alias
```

Plain text response (default):
```
files/
├── docs/
│   └── readme.txt (1.2 KB)
├── tools/
│   ├── nmap (856 KB)
│   └── nc (64 KB)
└── exploits/
    └── payload.py (2.1 KB)
```

JSON response (with `?format=json` or `Accept: application/json`):
```json
{
  "success": true,
  "root": { ... },
  "tree_string": "files/\n├── docs/\n│   └── readme.txt (1.2 KB)\n..."
}
```

### File Upload

Upload files via multipart form data:

```bash
POST /api/v1/upload
Content-Type: multipart/form-data

# Using curl
curl -X POST -F "file=@example.txt" http://localhost:8080/api/v1/upload
```

Response:
```json
{
  "success": true,
  "filename": "example.txt",
  "size": 1024,
  "path": "./uploads/example.txt"
}
```

### Uploads List

List all uploaded files (defaults to human-readable format):

```bash
GET /api/v1/uploads
GET /api/v1/ul        # Short alias
```

Plain text response (default):
```
Uploaded Files (3):
├── exploit.py (2.1 KB) - 2025-08-04 10:30:15
├── payload.sh (856 B) - 2025-08-04 10:25:42
└── data.zip (1.2 MB) - 2025-08-04 10:20:10
```

JSON response (with `?format=json` or `Accept: application/json`):
```json
{
  "success": true,
  "files": [
    {
      "name": "exploit.py",
      "size": 2148,
      "mod_time": "2025-08-04T10:30:15Z",
      "size_human": "2.1 KB"
    }
  ],
  "count": 3
}
```

### File Download

Download files via static file serving:

```bash
GET /files/path/to/file.txt
```

## Usage Examples

### Upload a file

```bash
# Upload a single file
curl -X POST -F "file=@payload.sh" http://localhost:8080/api/v1/upload

# Upload with verbose output
curl -v -X POST -F "file=@exploit.py" http://localhost:8080/api/v1/upload
```

### Download a file

```bash
# Download directly
curl -o downloaded_file.txt http://localhost:8080/files/example.txt

# Or use wget
wget http://localhost:8080/files/tools/nmap
```

### Browse directory structure

```bash
# Get file tree as JSON (structured data)
curl http://localhost:8080/api/v1/filetree | jq .

# Get human-readable tree (plain text, default)
curl http://localhost:8080/api/v1/tree

# Short Unix-style alias
curl http://localhost:8080/api/v1/ls

# List uploaded files (human-readable, default)
curl http://localhost:8080/api/v1/uploads

# Short alias for uploads
curl http://localhost:8080/api/v1/ul

# Get uploads as JSON
curl "http://localhost:8080/api/v1/uploads?format=json" | jq .

# Get health status
curl http://localhost:8080/api/v1/health
```

### Python upload script

```python
import requests

def upload_file(file_path, server_url="http://localhost:8080"):
    with open(file_path, 'rb') as f:
        files = {'file': f}
        response = requests.post(f"{server_url}/api/v1/upload", files=files)
        return response.json()

# Usage
result = upload_file("exploit.py")
print(f"Upload successful: {result['success']}")
```

## Security Considerations

- **File Validation**: Filenames are sanitized to prevent path traversal attacks
- **Size Limits**: Configurable upload size limits prevent DoS attacks
- **Directory Restrictions**: Uploads are contained within the designated upload directory
- **Input Sanitization**: All user inputs are properly validated

## Development

### Project Structure

```
ctfserver/
├── main.go                 # Application entry point
├── pkg/
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP request handlers
│   ├── logger/            # Logging and middleware
│   ├── models/            # Data structures
│   ├── server/            # HTTP server setup
│   ├── service/           # Business logic
│   └── util/              # Utility functions
└── README.md
```

### Adding New Features

1. **New API Endpoints**: Add handlers in `pkg/handlers/`
2. **Business Logic**: Add services in `pkg/service/`
3. **Data Models**: Define structures in `pkg/models/`
4. **Configuration**: Update `pkg/config/config.go`

## License

This project is designed for educational and CTF purposes. Use responsibly.e server written in Go

## Usage

- List all files

```bash
curl -O http://localhost:8080/filetree
```

- Upload a file

```bash
curl -X POST -F "file=@/path/to/file" http://localhost:8080/upload
```

# TODO

- [ ] Add HTTPS
