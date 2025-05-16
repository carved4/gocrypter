#!/bin/bash

echo "Building C++ RunPE library..."

# Navigate to cpp directory
cd cpp || exit 1

# Check if we're on Windows
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || "$OSTYPE" == "cygwin" ]]; then
    # Windows with MinGW or similar
    echo "Detected Windows environment"
    
    # Compile the C++ file into an object file
    g++ -c runPE.cpp -o runpe.o -Wall -std=c++11
    
    # Create static library
    ar rcs librunpe.a runpe.o
    
    # Copy the library to the parent directory for Go to find it
    cp librunpe.a ..
    
    # Clean up temporary object file
    rm -f runpe.o
else
    echo "Error: This library is only supported on Windows platforms"
    exit 1
fi

echo "C++ library build completed"
exit 0 