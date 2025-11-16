#!/usr/bin/env bash

# Script to run the homomorphic encryption variance calculation example
# Project: PrivacyLab

set -e

echo "========================================================"
echo "  Variance Calculation with Homomorphic Encryption"
echo "  Using SIMD Batching for Enhanced Privacy"
echo "========================================================"
echo ""

# Script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Verify Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    echo "Please install Go from https://go.dev/dl/"
    exit 1
fi

echo "Go detected: $(go version)"
echo ""

# Navigate to project root
cd "$SCRIPT_DIR/../.."

# Install or update dependencies
echo "Installing dependencies (Lattigo v6)..."
go get github.com/tuneinsight/lattigo/v6
go mod tidy
echo ""

# Return to example directory
cd "$SCRIPT_DIR"

# Generate population file if it doesn't exist or regenerate it
echo "Generating population data..."
cat > population.txt << EOF
10.5
15.2
20.7
25.3
30.1
35.8
40.2
45.6
50.9
55.4
EOF
echo "Population file created with 10 values"
echo ""

# Compile
echo "Compiling example..."
go build -o variance-example main.go
echo ""

# Run examples with different input values
echo "========================================================"
echo "Running Example 1: x = 5.0"
echo "========================================================"
echo ""
echo "5.0" | ./variance-example
echo ""

echo "========================================================"
echo "Running Example 2: x = 30.0"
echo "========================================================"
echo ""
echo "30.0" | ./variance-example
echo ""

echo "========================================================"
echo "Running Example 3: x = 50.0"
echo "========================================================"
echo ""
echo "50.0" | ./variance-example
echo ""

# Cleanup
echo ""
echo "========================================================"
echo "All examples completed successfully"
echo "========================================================"
echo ""
echo "Key features demonstrated:"
echo "  - All population values encrypted in 1 ciphertext (SIMD batching)"
echo "  - Population size hidden from observers"
echo "  - Variance calculated without revealing original data"
echo "  - Results verified against plain-text calculations"
echo ""
