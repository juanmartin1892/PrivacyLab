# Operation Domain# Operation Domain



This package defines statistical operations that can be performed on datasets in plaintext. The domain layer focuses purely on the mathematical logic without infrastructure concerns like cryptographic requirements.This package defines statistical operations that can be performed on datasets, both in plaintext and using homomorphic encryption.



## Purpose## Purpose



Provides abstractions for statistical computations that:Provides abstractions for statistical computations that:

- Define a common interface for all statistical operations- Define a common interface for all statistical operations

- Support plaintext computation for verification and testing- Support both plaintext and encrypted computation

- Validate operation parameters before execution- Specify required cryptographic resources (rotation keys)

- Remain independent of infrastructure concerns- Validate operation parameters before execution



## Key Concepts## Key Concepts



### StatisticalOperation Interface### StatisticalOperation Interface

Defines the contract that all statistical operations must implement:Defines the contract that all statistical operations must implement:

- **Type()**: Identifies the operation using a string-based identifier- **Type()**: Identifies the operation (variance, mean, etc.)

- **ComputePlaintext()**: Executes the operation on unencrypted data- **ComputePlaintext()**: Executes the operation on unencrypted data

- **Validate()**: Checks if operation can be performed with given parameters- **RequiredRotations()**: Specifies rotation keys needed for encrypted computation

- **Validate()**: Checks if operation can be performed with given parameters

### OperationType

String-based identifier that allows dynamic registration of operations without modifying the domain layer. New operations can define their own type constants without changing core domain code.### OperationParams

Encapsulates parameters specific to each operation:

### OperationParams- **Data**: Flexible map for operation-specific parameters

Encapsulates parameters specific to each operation:

- **Data**: Flexible map for operation-specific parameters### Implemented Operations



### Implemented Operations#### Variance

Calculates the variance of a dataset against a reference value.

#### Variance

Calculates the variance of a dataset against a reference value.Formula: `Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)`



Formula: `Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)`Requirements:

- Reference value x provided as first element in the data array

Type identifier: `"variance"` (defined as `operation.VarianceType`)- Rotation keys for positions 1 to n-1 (for SIMD slot summation)



Requirements:## Design Decisions

- Reference value x provided as first element in the data array

- Operations are stateless and can be reused

## Design Decisions- Plaintext computation provided for verification and testing

- Rotation requirements exposed before encryption to optimize key generation

### Separation of Concerns- Validation separated from computation for early error detection

The domain layer is now completely separated from infrastructure concerns:- Interface design allows easy extension for new operations (mean, standard deviation, etc.)

- No cryptographic concepts in domain operations

- No knowledge of rotation keys or SIMD operations## Usage

- Pure mathematical operations only

```go

Infrastructure adapters (like `lattigo.RotationMapper`) handle cryptographic requirements separately.import "github.com/juanmartin/privacylab/internal/domain/operation"



### Dynamic Registration// Create operation

Operations use string-based type identifiers instead of hardcoded constants, allowing:variance := operation.NewVarianceOperation()

- New operations to be added without modifying domain code

- External packages to define custom operation types// Compute in plaintext (reference value as first element, then population)

- Runtime registration of operations through the registry serviceparams := operation.NewOperationParams()

dataWithReference := append([]float64{5.0}, []float64{1, 2, 3, 4, 5}...)

### Stateless Designresult, err := variance.ComputePlaintext(dataWithReference, params)

- Operations are stateless and can be reused

- Plaintext computation provided for verification and testing// Get required rotations for encrypted computation

- Validation separated from computation for early error detectionrotations := variance.RequiredRotations(5)

```

## Usage

```go
import "github.com/juanmartin/privacylab/internal/domain/operation"

// Create operation
variance := operation.NewVarianceOperation()

// Compute in plaintext (reference value as first element, then population)
params := operation.NewOperationParams()
dataWithReference := append([]float64{5.0}, []float64{1, 2, 3, 4, 5}...)
result, err := variance.ComputePlaintext(dataWithReference, params)

// Get operation type identifier
opType := variance.Type() // Returns "variance"
```

## Extending with New Operations

To add a new statistical operation:

1. Define a type constant in your operation file:
```go
const MeanType OperationType = "mean"
```

2. Implement the `StatisticalOperation` interface:
```go
type MeanOperation struct{}

func (m *MeanOperation) Type() OperationType {
    return MeanType
}

func (m *MeanOperation) ComputePlaintext(data []float64, params OperationParams) (float64, error) {
    // Implementation
}

func (m *MeanOperation) Validate(dataSize int, params OperationParams) error {
    // Validation logic
}
```

3. Register with the registry service:
```go
registry.Register(MeanType, NewMeanOperation())
```

4. If your operation requires specific cryptographic resources, add a strategy to `lattigo.RotationMapper`:
```go
rotationMapper.Register(MeanType, func(dataSize int) []int {
    // Return required rotation positions
})
```
