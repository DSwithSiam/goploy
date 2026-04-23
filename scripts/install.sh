#!/bin/sh
set -e

REPO="DSwithSiam/goploy"
BINARY="goploy"

# Detect OS
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
  ARCH=amd64
fi

if [ "$OS" = "darwin" ]; then
  FILE="goploy-darwin-amd64"
elif [ "$OS" = "linux" ]; then
  FILE="goploy-linux-amd64"
else
  echo "Unsupported OS: $OS"
  exit 1
fi

LATEST=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep browser_download_url | grep $FILE | cut -d '"' -f 4)
if [ -z "$LATEST" ]; then
  echo "Could not find release for $FILE"
  exit 1
fi

echo "Downloading $LATEST..."
curl -L "$LATEST" -o "$BINARY"
chmod +x "$BINARY"
sudo mv "$BINARY" /usr/local/bin/

if command -v goploy >/dev/null 2>&1; then
  echo "goploy installed successfully!"
  goploy --help
else
  echo "Install failed. Please check your permissions."
  exit 1
fi
