# Build stage
FROM golang:1.24.1-alpine AS builder

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

# Install ca-certificates for HTTPS and shadow for user management
RUN apk --no-cache add ca-certificates shadow

# Create a non-root user with specific UID/GID that matches host user
ARG USER_ID=1000
ARG GROUP_ID=1000
RUN addgroup -g ${GROUP_ID} ctfuser && \
  adduser -D -u ${USER_ID} -G ctfuser ctfuser

# Create directories with proper ownership
RUN mkdir -p /opt/tools /opt/loot && \
  chown -R ctfuser:ctfuser /opt/tools /opt/loot

# Copy the binary from builder stage
COPY --from=builder /app/ctfserver /usr/local/bin/ctfserver
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /usr/local/bin/ctfserver /entrypoint.sh

# Switch to non-root user
USER ctfuser
WORKDIR /home/ctfuser

# Expose port
EXPOSE 80

# Set environment variables
ENV CTF_HOST=0.0.0.0
ENV CTF_PORT=80
ENV CTF_ROOT_DIR=/opt/tools
ENV CTF_UPLOAD_DIR=/opt/loot
ENV CTF_LOG_LEVEL=info

# Use entrypoint script
ENTRYPOINT ["/entrypoint.sh"]

# Run the application
CMD ["ctfserver"]
