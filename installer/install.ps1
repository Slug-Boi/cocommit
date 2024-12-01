$binaryUrl = "https://github.com/Slug-Boi/cocommit/blob/main/installer/bin/install-win"
$outputPath = "install-win.exe"

# Download the binary
Invoke-WebRequest -Uri $binaryUrl -OutFile $outputPath

# Run the binary
& .\$outputPath