# Determine the OS and architecture

# Set up a cleanup function to be triggered upon script exit
function Cleanup {
    Remove-Item -ErrorAction SilentlyContinue "cocommit.tar.gz"
    Remove-Item -ErrorAction SilentlyContinue "author.txt"
    if ($file) {
        Remove-Item -ErrorAction SilentlyContinue $file
    }
}

trap { Cleanup } EXIT

$OS = (Get-CimInstance Win32_OperatingSystem).Caption
$ARCH = (Get-CimInstance Win32_Processor).Architecture

$file = ""

$url = "https://github.com/Slug-Boi/cocommit/releases/latest/download/"

# Set the download URL based on the OS and architecture
if ($OS -match "Windows") {
    $URL = "${url}cocommit-win.tar.gz"
    $file = "cocommit.exe"
} else {
    Write-Host "Unsupported OS: $OS"
    exit 1
}

# Download and run the script
Invoke-WebRequest -Uri $URL -OutFile "cocommit.tar.gz"
if ($?) {
    tar -xvzf "cocommit.tar.gz"
    Remove-Item "cocommit.tar.gz"
    Remove-Item "author.txt"
    if ($file) {
        & .\$file -v
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Failed to extract the binary"
            exit 1
        }
    }
} else {
    Write-Host "Failed to download the file"
    exit 1
}

# Move the binary to the specified directory
$target_dir = Read-Host "Enter the directory to move the binary to (default: C:\Program Files\cocommit)"
$target_dir = if ($target_dir) { $target_dir } else { "C:\Program Files\cocommit" }

if (-Not (Test-Path (Split-Path $target_dir))) {
    Write-Host "Directory does not exist: $(Split-Path $target_dir)"
    exit 1
}

Move-Item -Path $file -Destination $target_dir
if ($?) {
    Write-Host "Binary moved to $target_dir successfully"
} else {
    Write-Host "Failed to move the binary to $target_dir"
    exit 1
}