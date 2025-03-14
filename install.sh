
#!/bin/bash
set -e

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Get latest version from GitHub releases
VERSION=$(curl -s https://api.github.com/repos/inodinwetrust10/godex/releases/latest \
          | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
VERSION=${VERSION#v}  # Remove 'v' prefix if present

# Download the appropriate binary archive
BINARY="godex_${OS}_${ARCH}"
URL="https://github.com/inodinwetrust10/godex/releases/download/v${VERSION}/${BINARY}.tar.gz"

echo "Downloading godex v${VERSION} for ${OS}/${ARCH}..."
TEMP_FILE=$(mktemp)
curl -L "$URL" -o "$TEMP_FILE"

# Validate the downloaded file is a gzip archive
if ! file "$TEMP_FILE" | grep -q 'gzip compressed data'; then
  echo "Downloaded file is not a valid gzip archive. Please check the URL and asset name."
  rm "$TEMP_FILE"
  exit 1
fi

tar xzf "$TEMP_FILE"
rm "$TEMP_FILE"

# Check if a previous installation exists
if [ -f /usr/local/bin/godex ]; then
  echo "A previous installation of godex was found at /usr/local/bin/godex."
  read -p "Do you want to update it to v${VERSION}? [y/N] " answer
  if [[ "$answer" =~ ^[Yy]$ ]]; then
    echo "Removing previous installation..."
    sudo rm -f /usr/local/bin/godex
  else
    echo "Installation aborted."
    exit 0
  fi
fi

# Install the new version
echo "Installing godex to /usr/local/bin..."
sudo mv "${BINARY}" /usr/local/bin/godex
chmod +x /usr/local/bin/godex

echo "godex v${VERSION} has been installed!"

