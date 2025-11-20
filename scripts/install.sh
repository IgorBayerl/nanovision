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

# Fetch Release Info
LATEST_JSON=$(curl -s -H "User-Agent: nanovision-installer" https://api.github.com/repos/$REPO/releases/latest)

if [[ "$LATEST_JSON" == *"<!DOCTYPE html"* ]]; then
    echo -e "${RED}Error: GitHub API unavailable.${NC}"
    exit 1
fi

if echo "$LATEST_JSON" | grep -q '"message": "Not Found"'; then
    echo -e "${RED}Error: No public release found.${NC}"
    exit 1
fi

LATEST_VERSION=$(echo "$LATEST_JSON" | grep -o '"tag_name": *"[^"]*"' | head -n 1 | sed 's/"tag_name": *"//' | sed 's/"//')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Error: Could not parse version.${NC}"
    exit 1
fi

# Version Check
CURRENT_VERSION="none"
if command -v $BIN_NAME &> /dev/null; then
    CURRENT_VERSION=$($BIN_NAME --version | grep -o 'v[0-9.]*' | head -n 1)
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

# Architecture
ARCH=$(uname -m)
if [ "$ARCH" == "x86_64" ]; then
    ASSET_NAME="nanovision_${LATEST_VERSION}_linux_amd64.tar.gz"
elif [ "$ARCH" == "aarch64" ]; then
    ASSET_NAME="nanovision_${LATEST_VERSION}_linux_arm64.tar.gz"
else
    echo -e "${RED}Unsupported architecture: $ARCH${NC}"
    exit 1
fi

DOWNLOAD_URL=$(echo "$LATEST_JSON" | grep -o "\"browser_download_url\": *\"[^\"]*$ASSET_NAME\"" | head -n 1 | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo -e "${RED}Error: Download URL not found for $ASSET_NAME${NC}"
    exit 1
fi

echo -e "${BLUE}Downloading...${NC}"

# Create Temp Dir
TMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# Download
curl -L -H "User-Agent: nanovision-installer" -o "$TMP_DIR/download.tar.gz" "$DOWNLOAD_URL"

# Extract All
tar xf "$TMP_DIR/download.tar.gz" -C "$TMP_DIR"

# Find binary
FOUND_BIN=$(find "$TMP_DIR" -type f -name "$BIN_NAME" | head -n 1)

if [ -z "$FOUND_BIN" ]; then
    echo -e "${RED}Error: Binary '$BIN_NAME' not found in archive.${NC}"
    exit 1
fi

# Install
mkdir -p "$INSTALL_DIR"
mv "$FOUND_BIN" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

echo -e "${GREEN}Installed $LATEST_VERSION to $INSTALL_DIR${NC}"

# Setup PATH automatically
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    SHELL_CONFIG=""
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        SHELL_CONFIG="$HOME/.bashrc"
    else
        # Detect via shell path
        case "$SHELL" in
        */zsh) SHELL_CONFIG="$HOME/.zshrc" ;;
        */bash) SHELL_CONFIG="$HOME/.bashrc" ;;
        *) SHELL_CONFIG="$HOME/.profile" ;;
        esac
    fi

    echo -e "${BLUE}Adding $INSTALL_DIR to PATH in $SHELL_CONFIG...${NC}"
    
    # Append to config file safely
    echo "" >> "$SHELL_CONFIG"
    echo "# Nanovision CLI" >> "$SHELL_CONFIG"
    echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_CONFIG"
    
    echo -e "${GREEN}Success! To use nanovision immediately, run:${NC}"
    echo "  source $SHELL_CONFIG"
fi