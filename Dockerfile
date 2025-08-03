# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ctfserver .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/ctfserver .

# Create directories for uploads and files
RUN mkdir -p uploads files

# Expose port
EXPOSE 8080

# Set environment variables
ENV CTF_HOST=0.0.0.0
ENV CTF_PORT=8080
ENV CTF_ROOT_DIR=/root/files
ENV CTF_UPLOAD_DIR=/root/uploads
ENV CTF_LOG_LEVEL=info

# Run the application
CMD ["./ctfserver"]
