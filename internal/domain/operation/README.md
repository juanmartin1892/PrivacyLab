# Operation Domain# Operation Domain# Operation Domain



This package defines statistical operations that can be performed on datasets using homomorphic encryption.



## PurposeThis package defines statistical operations that can be performed on datasets in plaintext. The domain layer focuses purely on the mathematical logic without infrastructure concerns like cryptographic requirements.This package defines statistical operations that can be performed on datasets, both in plaintext and using homomorphic encryption.



Provides abstractions for statistical computations that:



- Define a common interface for all statistical operations## Purpose## Purpose

- Remain independent of implementation details

- Allow dynamic registration of new operations



## Key ConceptsProvides abstractions for statistical computations that:Provides abstractions for statistical computations that:



### StatisticalOperation Interface- Define a common interface for all statistical operations- Define a common interface for all statistical operations



Defines the minimal contract that all statistical operations must implement:- Support plaintext computation for verification and testing- Support both plaintext and encrypted computation



- **Type()**: Identifies the operation using a string-based identifier- Validate operation parameters before execution- Specify required cryptographic resources (rotation keys)



This interface focuses only on operation identification. Implementation details such as:- Remain independent of infrastructure concerns- Validate operation parameters before execution

- Plaintext computation (for verification)

- Validation logic

- Cryptographic operations

## Key Concepts## Key Concepts

Are delegated to infrastructure adapters where they belong.



### OperationType

### StatisticalOperation Interface### StatisticalOperation Interface

String-based identifier that allows dynamic registration of operations without modifying the domain layer. New operations can define their own type constants without changing core domain code.

Defines the contract that all statistical operations must implement:Defines the contract that all statistical operations must implement:

### OperationParams

- **Type()**: Identifies the operation using a string-based identifier- **Type()**: Identifies the operation (variance, mean, etc.)

Encapsulates parameters specific to each operation:

- **ComputePlaintext()**: Executes the operation on unencrypted data- **ComputePlaintext()**: Executes the operation on unencrypted data

- **Data**: Flexible map for operation-specific parameters

- **Validate()**: Checks if operation can be performed with given parameters- **RequiredRotations()**: Specifies rotation keys needed for encrypted computation

### Implemented Operations

- **Validate()**: Checks if operation can be performed with given parameters

#### Variance

### OperationType

Calculates the variance of a dataset against a reference value.

String-based identifier that allows dynamic registration of operations without modifying the domain layer. New operations can define their own type constants without changing core domain code.### OperationParams

Formula: `Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)`

Encapsulates parameters specific to each operation:

Type identifier: `"variance"` (defined as `operation.VarianceType`)

### OperationParams- **Data**: Flexible map for operation-specific parameters

## Design Decisions

Encapsulates parameters specific to each operation:

### Separation of Concerns

- **Data**: Flexible map for operation-specific parameters### Implemented Operations

The domain layer is now completely separated from infrastructure concerns:



- No plaintext computation in domain (belongs to examples/tests)

- No validation logic in domain (belongs to infrastructure adapters)### Implemented Operations#### Variance

- No cryptographic concepts in domain operations

- Pure type identification onlyCalculates the variance of a dataset against a reference value.



Infrastructure adapters handle:#### Variance

- Cryptographic requirements (rotation keys, SIMD operations)

- Validation logic specific to the implementationCalculates the variance of a dataset against a reference value.Formula: `Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)`

- Actual computation (both plaintext for testing and encrypted)



### Dynamic Registration

Formula: `Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)`Requirements:

Operations use string-based type identifiers instead of hardcoded constants, allowing:

- Reference value x provided as first element in the data array

- New operations to be added without modifying domain code

- External packages to define custom operation typesType identifier: `"variance"` (defined as `operation.VarianceType`)- Rotation keys for positions 1 to n-1 (for SIMD slot summation)

- Runtime registration of operations through the registry service



### Minimal Interface

Requirements:## Design Decisions

Operations are now minimal type identifiers:

- Stateless and lightweight- Reference value x provided as first element in the data array

- Easy to register and look up

- Implementation-agnostic- Operations are stateless and can be reused



## Usage## Design Decisions- Plaintext computation provided for verification and testing



```go- Rotation requirements exposed before encryption to optimize key generation

import "github.com/juanmartin/privacylab/internal/domain/operation"

### Separation of Concerns- Validation separated from computation for early error detection

// Create operation for registration

variance := operation.NewVarianceOperation()The domain layer is now completely separated from infrastructure concerns:- Interface design allows easy extension for new operations (mean, standard deviation, etc.)



// Get operation type identifier- No cryptographic concepts in domain operations

opType := variance.Type() // Returns "variance"

```- No knowledge of rotation keys or SIMD operations## Usage



## Extending with New Operations- Pure mathematical operations only



To add a new statistical operation:```go



1. Define a type constant and minimal implementation:Infrastructure adapters (like `lattigo.RotationMapper`) handle cryptographic requirements separately.import "github.com/juanmartin/privacylab/internal/domain/operation"

```go

const MeanType OperationType = "mean"



type meanOperation struct{}### Dynamic Registration// Create operation



func NewMeanOperation() *meanOperation {Operations use string-based type identifiers instead of hardcoded constants, allowing:variance := operation.NewVarianceOperation()

    return &meanOperation{}

}- New operations to be added without modifying domain code



func (m *meanOperation) Type() OperationType {- External packages to define custom operation types// Compute in plaintext (reference value as first element, then population)

    return MeanType

}- Runtime registration of operations through the registry serviceparams := operation.NewOperationParams()

```

dataWithReference := append([]float64{5.0}, []float64{1, 2, 3, 4, 5}...)

2. Register with the registry service:

```go### Stateless Designresult, err := variance.ComputePlaintext(dataWithReference, params)

registry.Register(MeanType, NewMeanOperation())

```- Operations are stateless and can be reused



3. Implement the actual computation logic in infrastructure adapters:- Plaintext computation provided for verification and testing// Get required rotations for encrypted computation

```go

// In lattigo.EvaluatorAdapter- Validation separated from computation for early error detectionrotations := variance.RequiredRotations(5)

func (e *EvaluatorAdapter) EvaluateOperation(op operation.StatisticalOperation, request *computation.Request) (*computation.Result, error) {

    switch op.Type() {```

    case operation.MeanType:

        return e.evaluateMean(request)## Usage

    // ...

    }```go

}import "github.com/juanmartin/privacylab/internal/domain/operation"

```

// Create operation

4. If needed, add plaintext verification in examples/tests:variance := operation.NewVarianceOperation()

```go

func computePlaintextMean(data []float64) float64 {// Compute in plaintext (reference value as first element, then population)

    sum := 0.0params := operation.NewOperationParams()

    for _, v := range data {dataWithReference := append([]float64{5.0}, []float64{1, 2, 3, 4, 5}...)

        sum += vresult, err := variance.ComputePlaintext(dataWithReference, params)

    }

    return sum / float64(len(data))// Get operation type identifier

}opType := variance.Type() // Returns "variance"

``````


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
