#!/bin/bash
set -e

REPO="IgorBayerl/nanovision"
INSTALL_DIR="$HOME/.local/bin"
BIN_NAME="nanovision"

echo "Checking for updates..."

LATEST_JSON=$(curl -s https://api.github.com/repos/$REPO/releases/latest)
LATEST_VERSION=$(echo "$LATEST_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Error: Could not find latest release."
    exit 1
fi

CURRENT_VERSION="none"
if command -v $BIN_NAME &> /dev/null; then
    # Assumes 'nanovision --version' returns something like "nanovision version v0.1.0"
    CURRENT_VERSION=$($BIN_NAME --version | grep -o 'v[0-9.]*')
fi

if [ "$CURRENT_VERSION" == "$LATEST_VERSION" ]; then
    echo "Already on latest version ($LATEST_VERSION)."
    exit 0
fi

if [ "$CURRENT_VERSION" != "none" ]; then
    echo "Update available: $CURRENT_VERSION -> $LATEST_VERSION"
    read -p "View Changelog? (y/N) " show_log
    if [[ $show_log =~ ^[Yy]$ ]]; then
        echo -e "\n--- Changelog ---"
        echo "$LATEST_JSON" | grep '"body":' | sed 's/\\r\\n/\n/g' | sed 's/\\"/"/g'
        echo -e "-----------------\n"
    fi
    read -p "Install Update? (Y/n) " confirm
    if [[ $confirm =~ ^[Nn]$ ]]; then exit 0; fi
fi

ASSET_NAME="nanovision_${LATEST_VERSION}_linux_amd64.tar.gz"
DOWNLOAD_URL=$(echo "$LATEST_JSON" | grep "browser_download_url" | grep "$ASSET_NAME" | head -n 1 | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Download URL not found for $ASSET_NAME"
    exit 1
fi

echo "Downloading..."
mkdir -p "$INSTALL_DIR"
# Download and extract, assuming the tar contains the binary at root level
curl -L "$DOWNLOAD_URL" | tar xz -C "$INSTALL_DIR" "$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"

echo "Installed $LATEST_VERSION to $INSTALL_DIR"

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "$INSTALL_DIR is not in your PATH."
    echo "Add 'export PATH=\$PATH:$INSTALL_DIR' to your .bashrc or .zshrc"
fi