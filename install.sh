#!/bin/bash
set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}┌────────────────────────────────────────┐${NC}"
    echo -e "${BLUE}│        ${GREEN}godex Installer ${VERSION}${BLUE}            │${NC}"
    echo -e "${BLUE}└────────────────────────────────────────┘${NC}"
    echo
}

print_step() {
    echo -e "${YELLOW}[${1}/${2}]${NC} ${3}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

print_step 1 6 "Detecting system..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  print_error "Unsupported architecture: $ARCH"
fi

print_success "Detected $OS/$ARCH"

# Get latest version from GitHub releases
print_step 2 6 "Checking for latest version..."
VERSION=$(curl -s https://api.github.com/repos/inodinwetrust10/godex/releases/latest \
          | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
VERSION=${VERSION#v}  # Remove 'v' prefix if present

if [ -z "$VERSION" ]; then
    print_error "Failed to retrieve version information"
fi

print_success "Latest version: v${VERSION}"

# Download the appropriate binary archive
print_step 3 6 "Downloading godex..."
BINARY="godex_${OS}_${ARCH}"
URL="https://github.com/inodinwetrust10/godex/releases/download/v${VERSION}/${BINARY}.tar.gz"
echo -e "From: ${BLUE}$URL${NC}"

TEMP_FILE=$(mktemp)
curl -L "$URL" -o "$TEMP_FILE" --progress-bar

if ! file "$TEMP_FILE" | grep -q 'gzip compressed data'; then
  print_error "Downloaded file is not a valid gzip archive. Please check the URL and asset name."
fi

print_success "Download complete"

# Extract the binary
print_step 4 6 "Extracting files..."
tar xzf "$TEMP_FILE"
rm "$TEMP_FILE"

print_success "Extraction complete"

print_step 5 6 "Installing godex..."
if [ -f /usr/local/bin/godex ]; then
  echo -e "${YELLOW}A previous installation of godex was found.${NC}"
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
sudo mv "${BINARY}" /usr/local/bin/godex
sudo chmod +x /usr/local/bin/godex

print_success "godex v${VERSION} has been installed to /usr/local/bin/godex"

# Setup shell completion
print_step 6 6 "Setting up shell completion..."

# Create config directory if it doesn't exist
mkdir -p ~/.config/godex

# Detect shell
SHELL_TYPE=$(basename "$SHELL")
case "$SHELL_TYPE" in
    bash)
        /usr/local/bin/godex completion bash > ~/.config/godex/godex.bash
        
        # Check if completion is already in .bashrc
        if ! grep -q "godex completion" ~/.bashrc; then
            echo "source ~/.config/godex/godex.bash" >> ~/.bashrc
            print_success "Bash completion installed. It will be activated in new shell sessions."
            echo -e "${YELLOW}Run 'source ~/.bashrc' to enable completion in the current session.${NC}"
        else
            print_success "Bash completion already configured in ~/.bashrc"
        fi
        ;;
    zsh)
        /usr/local/bin/godex completion zsh > ~/.config/godex/godex.zsh
        
        # Check if completion is already in .zshrc
        if ! grep -q "godex completion" ~/.zshrc; then
            echo "source ~/.config/godex/godex.zsh" >> ~/.zshrc
            print_success "Zsh completion installed. It will be activated in new shell sessions."
            echo -e "${YELLOW}Run 'source ~/.zshrc' to enable completion in the current session.${NC}"
        else
            print_success "Zsh completion already configured in ~/.zshrc"
        fi
        ;;
    fish)
        # Create fish completions directory if it doesn't exist
        mkdir -p ~/.config/fish/completions
        /usr/local/bin/godex completion fish > ~/.config/fish/completions/godex.fish
        print_success "Fish completion installed."
        ;;
    *)
        echo -e "${YELLOW}Shell completion for $SHELL_TYPE not set up automatically.${NC}"
        echo -e "You can manually set it up with: ${BLUE}godex completion --help${NC}"
        ;;
esac

# Installation complete
echo
echo -e "${GREEN}┌────────────────────────────────────────┐${NC}"
echo -e "${GREEN}│       Installation Complete! 🎉        │${NC}"
echo -e "${GREEN}└────────────────────────────────────────┘${NC}"
echo
echo -e "Run ${BLUE}godex --help${NC} to get started."
echo
echo -e "${YELLOW}Google Drive Integration:${NC}"
echo -e "For Google Drive backup functionality, you need to:"
echo -e "1. Create a Google Cloud project and enable Google Drive API"
echo -e "2. Create OAuth 2.0 credentials (Desktop app)"
echo -e "3. Download credentials.json and place it in ${BLUE}~/.config/godex/${NC}"
echo
echo -e "Thank you for installing godex! 🚀"
