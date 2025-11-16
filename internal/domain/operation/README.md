# Operation Domain

This package defines statistical operations that can be performed on datasets, both in plaintext and using homomorphic encryption.

## Purpose

Provides abstractions for statistical computations that:
- Define a common interface for all statistical operations
- Support both plaintext and encrypted computation
- Specify required cryptographic resources (rotation keys)
- Validate operation parameters before execution

## Key Concepts

### StatisticalOperation Interface
Defines the contract that all statistical operations must implement:
- **Type()**: Identifies the operation (variance, mean, etc.)
- **ComputePlaintext()**: Executes the operation on unencrypted data
- **RequiredRotations()**: Specifies rotation keys needed for encrypted computation
- **Validate()**: Checks if operation can be performed with given parameters

### OperationParams
Encapsulates parameters specific to each operation:
- **ReferenceValue**: Main parameter for the operation (e.g., the value x in variance calculation)
- **AdditionalData**: Flexible map for operation-specific parameters

### Implemented Operations

#### Variance
Calculates the variance of a dataset against a reference value.

Formula: `Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)`

Requirements:
- Reference value x
- Rotation keys for positions 1 to n-1 (for SIMD slot summation)

## Design Decisions

- Operations are stateless and can be reused
- Plaintext computation provided for verification and testing
- Rotation requirements exposed before encryption to optimize key generation
- Validation separated from computation for early error detection
- Interface design allows easy extension for new operations (mean, standard deviation, etc.)

## Usage

```go
import "github.com/juanmartin/privacylab/internal/domain/operation"

// Create operation
variance := operation.NewVarianceOperation()

// Compute in plaintext
params := operation.NewOperationParams(5.0) // reference value x = 5.0
result, err := variance.ComputePlaintext([]float64{1, 2, 3, 4, 5}, params)

// Get required rotations for encrypted computation
rotations := variance.RequiredRotations(5)
```
