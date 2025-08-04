#!/bin/bash

# CTF Server Docker Setup Script
# This script creates the necessary directories with proper permissions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Setting up CTF Server directories...${NC}"

# Get current user ID and group ID
USER_ID=$(id -u)
GROUP_ID=$(id -g)

echo -e "${YELLOW}📋 User ID: $USER_ID, Group ID: $GROUP_ID${NC}"

# Create directories if they don't exist
echo -e "${YELLOW}📁 Creating directories...${NC}"
sudo mkdir -p /opt/tools /opt/loot

# Set ownership to current user
echo -e "${YELLOW}🔐 Setting ownership...${NC}"
sudo chown -R $USER_ID:$GROUP_ID /opt/tools /opt/loot

# Set proper permissions
echo -e "${YELLOW}🔑 Setting permissions...${NC}"
sudo chmod -R 755 /opt/tools /opt/loot

# Export environment variables for docker-compose
export UID=$USER_ID
export GID=$GROUP_ID

echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo -e "${GREEN}🐳 You can now start the container with:${NC}"
echo "   UID=$USER_ID GID=$GROUP_ID docker-compose up -d"
echo ""
echo -e "${GREEN}📂 Directories created:${NC}"
echo "   /opt/tools  - For files to download"
echo "   /opt/loot   - For uploaded files"
echo ""
echo -e "${GREEN}🌐 The server will be available at:${NC}"
echo "   http://localhost"
echo ""
echo -e "${YELLOW}💡 To add files for download, copy them to:${NC}"
echo "   /opt/tools/"
echo ""
echo -e "${YELLOW}📤 Uploaded files will appear in:${NC}"
echo "   /opt/loot/"

# Create directories if they don't exist
echo -e "${YELLOW}📁 Creating directories...${NC}"
sudo mkdir -p /opt/tools /opt/loot

# Set ownership to current user
echo -e "${YELLOW}🔧 Setting ownership...${NC}"
sudo chown -R $USER_ID:$GROUP_ID /opt/tools /opt/loot

# Set permissions
echo -e "${YELLOW}🔒 Setting permissions...${NC}"
sudo chmod -R 755 /opt/tools
sudo chmod -R 755 /opt/loot

# Export environment variables for docker-compose
export UID=$USER_ID
export GID=$GROUP_ID

echo -e "${GREEN}✅ Setup complete!${NC}"
echo -e "${YELLOW}📖 Directory structure:${NC}"
echo -e "  📂 /opt/tools - For tools and files to download (read-only in container)"
echo -e "  📂 /opt/loot  - For uploaded files and loot (read-write in container)"
echo
echo -e "${YELLOW}🐳 You can now run:${NC}"
echo -e "  ${GREEN}docker-compose up -d${NC}"
echo
echo -e "${YELLOW}🌐 Server will be available at:${NC}"
echo -e "  ${GREEN}http://localhost:8080${NC}"
echo
echo -e "${YELLOW}📚 Useful endpoints:${NC}"
echo -e "  • Health: ${GREEN}curl http://localhost:8080/api/v1/health${NC}"
echo -e "  • List tools: ${GREEN}curl http://localhost:8080/api/v1/ls${NC}"
echo -e "  • List loot: ${GREEN}curl http://localhost:8080/api/v1/ul${NC}"
echo -e "  • Upload: ${GREEN}curl -X POST -F \"file=@filename\" http://localhost:8080/api/v1/upload${NC}"
