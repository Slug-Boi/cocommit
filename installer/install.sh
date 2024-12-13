#!/usr/bin/env bash
# Determine the OS and architecture

# Set up a cleanup function to be triggered upon script exit
__cleanup ()
{
    rm "cocommit.tar.gz" 2>/dev/null
    rm author.txt 2>/dev/null
    if [ -n "$file" ]; then
        rm $file 2>/dev/null
    fi
}

trap __cleanup EXIT

OS=$(uname -s)
ARCH=$(uname -m)

file=""

url="https://github.com/Slug-Boi/cocommit/releases/latest/download/"

# Set the download URL based on the OS and architecture
if [ "$OS" == "Linux" ]; then
    URL="${url}cocommit-linux.tar.gz"
    file="cocommit-linux"
elif [ "$OS" == "Darwin" ]; then
    if [ "$ARCH" == "x86_64" ]; then
        URL="${url}cocommit-darwin-x86_64.tar.gz"
        file="cocommit-darwin"
    else
        URL="${url}cocommit-darwin-aarch64.tar.gz"
        file="cocommit-darwin-aarch64"
    fi
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# Download and run the script

curl -L -o cocommit.tar.gz $URL && \
tar -xvzf cocommit.tar.gz && \
rm cocommit.tar.gz && rm author.txt && \
chmod +x $file && ./$file -v
if [ $? -ne 0 ]; then
    echo "Failed to extract the binary"
    exit 1
fi

# Move the binary to the current directory
read -p "Enter the directory to move the binary to (default: /usr/local/bin/cocommit): " target_dir
target_dir=${target_dir:-/usr/local/bin/cocommit}

if [ ! -d "$(dirname "$target_dir")" ]; then
    echo "Directory does not exist: $(dirname "$target_dir")"
    exit 1
fi

sudo mv $file "$target_dir"
if [ $? -ne 0 ]; then
    echo "Failed to move the binary to $target_dir"
    exit 1
fi

echo "Binary moved to $target_dir successfully"

