#!/bin/bash

# CTF Server Test Script

SERVER_URL="http://localhost:8080"
API_URL="$SERVER_URL/api/v1"

echo "=== CTF Server API Test ==="
echo "Server URL: $SERVER_URL"
echo

# Function to check if server is running
check_server() {
    if curl -s "$API_URL/health" > /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Function to test health endpoint
test_health() {
    echo "1. Testing health endpoint..."
    response=$(curl -s "$API_URL/health")
    echo "Response: $response"
    if echo "$response" | jq -e '.status == "healthy"' > /dev/null 2>&1; then
        echo "✅ Health check passed"
    else
        echo "❌ Health check failed"
    fi
    echo
}

# Function to test file tree endpoint
test_filetree() {
    echo "2. Testing file tree endpoint..."
    response=$(curl -s "$API_URL/filetree")
    echo "Response: $response"
    if echo "$response" | jq -e '.success == true' > /dev/null 2>&1; then
        echo "✅ File tree test passed"
    else
        echo "❌ File tree test failed"
    fi
    echo
}

# Function to test pretty file tree endpoint
test_pretty_filetree() {
    echo "3. Testing pretty file tree endpoint..."
    
    # Test default plain text response
    echo "   Default format (plain text):"
    text_response=$(curl -s "$API_URL/filetree/pretty")
    echo "$text_response" | head -5
    if echo "$text_response" | grep -q "testfiles"; then
        echo "   ✅ Pretty file tree default (text) test passed"
    else
        echo "   ❌ Pretty file tree default (text) test failed"
    fi
    
    # Test JSON response with parameter
    echo "   JSON format (with ?format=json):"
    response=$(curl -s "$API_URL/filetree/pretty?format=json")
    if echo "$response" | jq -e '.success == true and .tree_string != null' > /dev/null 2>&1; then
        echo "   ✅ Pretty file tree JSON test passed"
    else
        echo "   ❌ Pretty file tree JSON test failed"
    fi
    
    # Test shorter endpoints
    echo "   Short aliases:"
    tree_response=$(curl -s "$API_URL/tree")
    ls_response=$(curl -s "$API_URL/ls")
    if echo "$tree_response" | grep -q "testfiles" && echo "$ls_response" | grep -q "testfiles"; then
        echo "   ✅ Short aliases (/tree, /ls) test passed"
    else
        echo "   ❌ Short aliases test failed"
    fi
    echo
}

# Function to test upload endpoint
test_upload() {
    echo "4. Testing upload endpoint..."
    
    # Create a test file
    test_file="/tmp/ctf_test_upload.txt"
    echo "This is a test upload file" > "$test_file"
    
    response=$(curl -s -X POST -F "file=@$test_file" "$API_URL/upload")
    echo "Response: $response"
    
    if echo "$response" | jq -e '.success == true' > /dev/null 2>&1; then
        echo "✅ Upload test passed"
    else
        echo "❌ Upload test failed"
    fi
    
    # Clean up
    rm -f "$test_file"
    echo
}

# Function to test uploads list endpoint
test_uploads_list() {
    echo "5. Testing uploads list endpoint..."
    
    # Test default plain text response
    echo "   Default format (plain text):"
    text_response=$(curl -s "$API_URL/uploads")
    echo "$text_response" | head -5
    if echo "$text_response" | grep -q -E "(Uploaded Files|No uploaded files found)"; then
        echo "   ✅ Uploads list default (text) test passed"
    else
        echo "   ❌ Uploads list default (text) test failed"
    fi
    
    # Test JSON response with parameter
    echo "   JSON format (with ?format=json):"
    response=$(curl -s "$API_URL/uploads?format=json")
    if echo "$response" | jq -e '.success == true and .count >= 0' > /dev/null 2>&1; then
        echo "   ✅ Uploads list JSON test passed"
        
        # Show uploaded files count
        count=$(echo "$response" | jq -r '.count')
        echo "   Found $count uploaded file(s)"
        
        # Show file names if any
        if [ "$count" -gt 0 ]; then
            echo "   Files:"
            echo "$response" | jq -r '.files[].name' | sed 's/^/     - /'
        fi
    else
        echo "   ❌ Uploads list JSON test failed"
    fi
    echo
}

# Function to test file download
test_download() {
    echo "6. Testing file download..."
    
    # Try to download a file from the test directory
    if curl -s "$SERVER_URL/files/readme.txt" | grep -q "This is a test file"; then
        echo "✅ Download test passed"
    else
        echo "❌ Download test failed"
    fi
    echo
}

# Main execution
echo "Checking if server is running..."
if check_server; then
    echo "✅ Server is running"
    echo
    
    test_health
    test_filetree
    test_pretty_filetree
    test_upload
    test_uploads_list
    test_download
    
    echo "=== Test completed ==="
else
    echo "❌ Server is not running. Please start the server first:"
    echo "   ./ctfserver -port 8080 -root ./testfiles -upload-dir ./uploads"
    exit 1
fi
