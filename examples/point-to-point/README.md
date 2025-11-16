# Variance Calculation with Homomorphic Encryption

This example demonstrates privacy-preserving variance calculation using the CKKS homomorphic encryption scheme. The program computes the variance of an input value `x` against a population `X` without revealing the original data.

## Overview

The program calculates the variance using the formula:

```
Var(X, x) = E[(X - x)²] = (1/n) × Σ(X_i - x)²
```

where:
- `X = [X_1, X_2, ..., X_n]` is the population (loaded from a file)
- `x` is the input value
- `n` is the population size

## Privacy Features

This implementation uses **SIMD (Single Instruction Multiple Data) batching** to enhance privacy:

1. **Single ciphertext for entire population**: All population values are packed into the slots of a single ciphertext, rather than encrypting each value individually.
2. **Encrypted population size**: The number of individuals `n` is also encrypted.
3. **Encrypted input**: The input value `x` is encrypted.

This approach prevents attackers from determining the population size by counting ciphertexts, significantly improving privacy protection.

## How It Works

1. **Load Population**: Read population values from `population.txt`
2. **Get Input**: Accept input value `x` from the user
3. **Plain-text Calculation**: Compute variance in plain text for verification
4. **Key Generation**: Generate cryptographic keys including rotation keys for SIMD operations
5. **Encryption with SIMD Batching**:
   - Pack all population values into a single ciphertext using available slots
   - Encrypt `n` (population size) as a separate ciphertext
   - Encrypt `x` as a ciphertext replicated across all slots
6. **Homomorphic Computation**:
   - Compute `(X - x)` for all values simultaneously (element-wise subtraction)
   - Square the differences: `(X - x)²` (element-wise multiplication)
   - Sum all squared differences using rotation operations
   - Divide by `n` to get the variance
7. **Decryption**: Decrypt the result to obtain the variance
8. **Verification**: Compare with plain-text result to verify correctness

## Mathematical Operations

The homomorphic computation performs the following operations on encrypted data:

```
encryptedDiff = encryptedPopulation - encryptedX          // Element-wise subtraction
encryptedSquared = encryptedDiff × encryptedDiff          // Element-wise multiplication
encryptedSum = RotateAndSum(encryptedSquared, n)          // Sum using rotations
encryptedVariance = encryptedSum × (1/n)                  // Scalar multiplication
```

All operations are performed on ciphertexts without decryption, preserving privacy throughout the computation.

## Files

- `main.go`: Complete implementation of variance calculation with homomorphic encryption
- `population.txt`: Population data (one value per line)
- `run.sh`: Script to build and run the example
- `configs/`: Configuration directory (for future use)

Note: The binary `variance-example` is not included in the repository and must be built locally.

## How to Execute

### Prerequisites

Before running this example, ensure you have:

- **Go 1.21 or later** installed on your system
- Access to the PrivacyLab repository root directory
- Internet connection for downloading dependencies (first run only)

### Execution Methods

#### Method 1: Using the automated script (recommended)

The `run.sh` script handles all setup, compilation, and execution automatically:

```bash
cd examples/point-to-point
chmod +x run.sh
./run.sh
```

This script will:
1. Verify Go installation
2. Install required dependencies (Lattigo v6)
3. Generate or verify the population data file
4. Compile the example binary
5. Run three demonstration cases with different input values (x = 5.0, 30.0, 50.0)
6. Display results and verification for each case

#### Method 2: Manual execution

For manual control or custom input values:

**Step 1: Navigate to the example directory**
```bash
cd examples/point-to-point
```

**Step 2: Install dependencies** (first time only)
```bash
cd ../..
go get github.com/tuneinsight/lattigo/v6
go mod tidy
cd examples/point-to-point
```

**Step 3: Ensure population data exists**
```bash
# Verify population.txt exists, or create it:
cat > population.txt << 'EOF'
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
```

**Step 4: Build the binary**
```bash
go build -o variance-example main.go
```

**Step 5: Run with interactive input**
```bash
./variance-example
```

When prompted, enter a numeric value:
```
Enter input value x: 30.0
```

**Step 6: Run with piped input** (non-interactive)
```bash
echo "30.0" | ./variance-example
```

### Quick Start

For the fastest execution with default settings:

```bash
cd examples/point-to-point && ./run.sh
```

### Troubleshooting

**Issue: Permission denied when running `run.sh`**
```bash
chmod +x run.sh
```

**Issue: Go command not found**

Install Go from https://go.dev/dl/ and ensure it is in your PATH.

**Issue: Module dependencies not found**

Run from the repository root:
```bash
go mod download
go mod tidy
```

**Issue: Population file missing**

The `run.sh` script creates this automatically. If running manually, create `population.txt` with one numeric value per line.

## Example Output

```
=== Variance Calculation with Homomorphic Encryption ===
Computing variance of input x against population X using encrypted data

Reading population from file...
Population loaded: [10.5 15.2 20.7 25.3 30.1 35.8 40.2 45.6 50.9 55.4] (size: 10)

Enter input value x: 30.0
Input x: 30.00

Calculating variance in plain text (for verification)...
Plain text variance: 217.089000

Generating cryptographic keys...
Generating rotation keys for SIMD operations...

Encrypting population data using SIMD batching...
Population encrypted successfully (all 10 values in 1 ciphertext)
Encrypting population size n...
Population size encrypted successfully
Encrypting input x...
Input x encrypted successfully

Calculating variance on encrypted data...
Decrypting result...

Results:
Plain text variance:      217.089000
Homomorphic variance:     217.089000
Absolute error:           0.000000041
Relative error:           0.000000%

Example completed successfully
Variance was calculated without revealing the original data
```

## Security Considerations

- **Data Privacy**: The original population values and input are never revealed during computation
- **Population Size Privacy**: Using SIMD batching hides the exact number of individuals from observers who only see ciphertexts
- **Computational Privacy**: All variance calculations are performed on encrypted data
- **Verification**: Plain-text comparison ensures computational correctness while maintaining privacy guarantees

## Technical Details

- **Encryption Scheme**: CKKS (Cheon-Kim-Kim-Song) for approximate arithmetic on encrypted real numbers
- **SIMD Batching**: Encodes multiple values into polynomial slots for parallel operations
- **Rotation Keys**: Enable slot rotations needed to sum values across the ciphertext
- **Parameters**:
  - LogN: 14 (16384 slots available)
  - Precision: 40-bit default scale
  - Multiplicative Depth: 5 levels

## Dependencies

This example uses the Lattigo library for homomorphic encryption:
```
github.com/tuneinsight/lattigo/v6
```

## Applications

This technique is useful in scenarios such as:

- **Medical data analysis** without exposing sensitive patient information
- **Cloud computing** where the provider is not fully trusted
- **Financial analysis** preserving transaction privacy
- **Privacy-preserving machine learning**

## Next Steps

Possible extensions:

- Calculate other statistics (standard deviation, covariance)
- Implement homomorphic linear regression
- Explore multi-party computation scenarios
- Integrate zero-knowledge proofs for verifiability
