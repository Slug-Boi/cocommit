#!/usr/bin/env bash
# Determine the OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

file=""

url="https://github.com/Slug-Boi/cocommit/raw/refs/heads/chore_install_script/installer/bin/"

# Set the download URL based on the OS and architecture
if [ "$OS" == "Linux" ]; then
    URL="${url}install-linux"
    file="install-linux"
elif [ "$OS" == "Darwin" ]; then
    if [ "$ARCH" == "x86_64" ]; then
        URL="${url}install-darwin-x86_64"
        file="install-darwin-x86_64"
    else
        URL="${url}install-darwin-aarch64"
        file="install-darwin-aarch64"
    fi
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# Download and run the script

echo $file
curl -LJO $URL && chmod +x $file && ./$file
