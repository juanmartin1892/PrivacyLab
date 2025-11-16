# Computation Domain

This package defines the request-response model for homomorphic computations between client and worker components.

## Purpose

Provides domain entities for:
- **Request**: Encapsulates computation requests sent from client to worker
- **Result**: Encapsulates computation results returned from worker to client
- **ResultMetadata**: Tracks computational cost and performance metrics

## Key Concepts

### Request
Represents a computation request that the worker must process without access to plaintext data.

Contains:
- **OperationType**: The statistical operation to perform (variance, mean, etc.)
- **EncryptedInputs**: All encrypted data needed for the computation
- **Parameters**: Operation-specific configuration
- **RequestID**: Unique identifier for request tracking and result correlation

The request object is designed to be serializable for transmission over the network.

### Result
Represents the outcome of a homomorphic computation performed by the worker.

Contains:
- **EncryptedResult**: The computed result in encrypted form
- **OperationType**: The operation that was performed
- **RequestID**: Correlation with the original request
- **Metadata**: Detailed information about the computation

### ResultMetadata
Tracks important metrics about the homomorphic computation:
- **ComputationTime**: Duration of the computation
- **MultiplicationsUsed**: Number of homomorphic multiplications performed
- **RotationsUsed**: Number of SIMD slot rotations performed
- **LevelsConsumed**: Multiplication depth consumed
- **FinalLevel**: Remaining multiplication depth
- **FinalScale**: Final scaling factor of the result
- **WorkerID**: Identifier of the worker that performed the computation
- **Timestamp**: When the computation was completed

This metadata is crucial for:
- Performance monitoring and optimization
- Debugging cryptographic parameter issues
- Capacity planning for worker resources
- Audit trails

## Design Decisions

- Request and Result are immutable after creation
- All encrypted data is copied defensively
- Metadata automatically captures timestamp
- RequestID enables correlation in distributed systems
- Parameters stored as generic map for flexibility
- Worker must never have access to plaintext data

## Usage

```go
import "github.com/juanmartin/privacylab/internal/domain/computation"

// Create a computation request
request := computation.NewRequest(
    operation.OperationVariance,
    []*crypto.EncryptedData{encPop, encX, encN},
    "req-123",
)
request.SetParameter("populationSize", 100)

// Create a result
metadata := computation.NewResultMetadata(
    time.Second * 2,  // computation time
    3,                 // multiplications
    99,                // rotations
    2,                 // levels consumed
    3,                 // final level
    1048576.0,         // final scale
    "worker-01",       // worker ID
)
result := computation.NewResult(
    encryptedVariance,
    operation.OperationVariance,
    "req-123",
    metadata,
)
```
