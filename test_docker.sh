#!/bin/bash

# Test Docker setup
echo "🧪 Testing Docker setup..."

# Build the container
echo "🔨 Building container..."
UID=$(id -u) GID=$(id -g) docker-compose build

# Start the container
echo "🚀 Starting container..."
UID=$(id -u) GID=$(id -g) docker-compose up -d

# Wait for container to start
echo "⏳ Waiting for server to start..."
sleep 5

# Test health endpoint
echo "🏥 Testing health endpoint..."
response=$(curl -s http://localhost/api/v1/health 2>/dev/null || echo "failed")

if echo "$response" | grep -q "healthy"; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed: $response"
    echo "📋 Container logs:"
    docker-compose logs ctfserver
    exit 1
fi

# Test file upload
echo "📤 Testing file upload..."
echo "test data" > /tmp/test_upload.txt
upload_response=$(curl -s -X POST -F "file=@/tmp/test_upload.txt" http://localhost/api/v1/upload 2>/dev/null || echo "failed")

if echo "$upload_response" | grep -q "success.*true"; then
    echo "✅ Upload test passed"
    
    # Check if file exists with correct permissions
    if [ -f "/opt/loot/test_upload.txt" ]; then
        owner=$(stat -c "%U:%G" /opt/loot/test_upload.txt)
        current_user=$(whoami)
        if [ "$owner" = "$current_user:$current_user" ]; then
            echo "✅ File permissions correct: $owner"
        else
            echo "⚠️  File permissions: $owner (expected: $current_user:$current_user)"
        fi
    else
        echo "❌ Uploaded file not found in /opt/loot/"
    fi
else
    echo "❌ Upload test failed: $upload_response"
fi

# Test uploads list
echo "📋 Testing uploads list..."
uploads_response=$(curl -s http://localhost/api/v1/uploads 2>/dev/null || echo "failed")

if echo "$uploads_response" | grep -q "test_upload.txt"; then
    echo "✅ Uploads list test passed"
else
    echo "❌ Uploads list test failed"
fi

# Cleanup
echo "🧹 Cleaning up..."
rm -f /tmp/test_upload.txt
docker-compose down

echo "🎉 Docker test completed!"
