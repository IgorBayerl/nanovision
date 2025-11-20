#!/bin/bash
set -e

REPO="IgorBayerl/nanovision"
INSTALL_DIR="$HOME/.local/bin"
BIN_NAME="nanovision"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}Checking for latest nanovision version...${NC}"

# User-Agent is required by GitHub API
LATEST_JSON=$(curl -s -H "User-Agent: nanovision-installer" https://api.github.com/repos/$REPO/releases/latest)

# Check if release exists (API returns "Not Found" for drafts or bad repos)
if echo "$LATEST_JSON" | grep -q '"message": "Not Found"'; then
    echo -e "${RED}Error: No public release found.${NC}"
    echo "The latest release might still be a Draft. Publish it on GitHub to make it available."
    exit 1
fi

LATEST_VERSION=$(echo "$LATEST_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Error: Could not parse version from GitHub response.${NC}"
    exit 1
fi

CURRENT_VERSION="none"
if command -v $BIN_NAME &> /dev/null; then
    CURRENT_VERSION=$($BIN_NAME --version | grep -o 'v[0-9.]*')
fi

if [ "$CURRENT_VERSION" == "$LATEST_VERSION" ]; then
    echo -e "${GREEN}You are already on the latest version ($LATEST_VERSION).${NC}"
    exit 0
fi

if [ "$CURRENT_VERSION" != "none" ]; then
    echo -e "Update available: $CURRENT_VERSION -> $LATEST_VERSION"
    read -p "Install Update? (Y/n) " confirm
    if [[ $confirm =~ ^[Nn]$ ]]; then exit 0; fi
fi

ARCH=$(uname -m)
if [ "$ARCH" == "x86_64" ]; then
    ASSET_NAME="nanovision_${LATEST_VERSION}_linux_amd64.tar.gz"
elif [ "$ARCH" == "aarch64" ]; then
    ASSET_NAME="nanovision_${LATEST_VERSION}_linux_arm64.tar.gz"
else
    echo -e "${RED}Unsupported architecture: $ARCH${NC}"
    exit 1
fi

DOWNLOAD_URL=$(echo "$LATEST_JSON" | grep "browser_download_url" | grep "$ASSET_NAME" | head -n 1 | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo -e "${RED}Error: Download URL not found for $ASSET_NAME${NC}"
    exit 1
fi

echo -e "${BLUE}Downloading...${NC}"
mkdir -p "$INSTALL_DIR"

curl -L -H "User-Agent: nanovision-installer" "$DOWNLOAD_URL" | tar xz -C "$INSTALL_DIR" "$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

echo -e "${GREEN}Installed $LATEST_VERSION to $INSTALL_DIR${NC}"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${BLUE}Note: $INSTALL_DIR is not in your PATH.${NC}"
    echo "Run this to add it:"
    echo "  export PATH=\$PATH:$INSTALL_DIR"
fi