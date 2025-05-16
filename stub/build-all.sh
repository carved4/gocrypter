#!/bin/bash

echo "=== RunPE Crypter Builder ==="

# Make sure the script is executable
chmod +x build.sh

# Build the C++ library first
echo "Building C++ library..."
./build.sh
if [ $? -ne 0 ]; then
    echo "Failed to build C++ library."
    exit 1
fi

# Check if we're on Windows (required for RunPE)
if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "win32" && "$OSTYPE" != "cygwin" ]]; then
    echo "Warning: This crypter is designed for Windows only."
    echo "You can build it on a non-Windows platform, but it will only run on Windows."
fi

# Now build the Go application
echo "Building Go application..."
go mod tidy

# Build for 64-bit Windows
echo "Building for Windows 64-bit..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -tags windows -ldflags="-s -w -H=windowsgui" -o runpe.exe
if [ $? -ne 0 ]; then
    echo "Failed to build Go application for amd64."
    exit 1
fi

echo "Build completed successfully!"
echo "Output file: runpe.exe (64-bit)"

# Cleanup temporary files
echo "Cleaning up temporary files..."
rm -f librunpe.a
rm -f key.txt encrypted_Input.bin 2>/dev/null

echo "Cleanup completed."

