# Process Service

The process service handles encrypted computation requests on the worker side. It coordinates between domain operation logic and homomorphic evaluation infrastructure.

## Purpose

Execute statistical operations on encrypted data without decrypting it, using homomorphic encryption capabilities.

## Responsibilities

- Validate that requested operations can be performed on provided data.
- Coordinate between operation definitions and homomorphic evaluator.
- Track computation performance (time, operations, memory).
- Delegate operation management to the registry service.
- Maintain worker identity for result traceability.

## Service API

### Constructor

```go
NewService(
    evaluator ports.HomomorphicEvaluator,
    workerID string,
    operationRegistry *registry.Service,
) *Service
```

Creates a new processing service with the required infrastructure dependencies.
If operationRegistry is nil, an empty registry is created automatically.

### Methods

- **ProcessRequest**: Processes an encrypted computation request and returns encrypted result.
- **SupportedOperations**: Returns list of currently supported operation types from the registry.
- **WorkerID**: Returns the identifier of this worker instance.

## Usage Example

```go
// Create registry with default operations
operationRegistry := registry.NewService().RegisterDefaults()

// Create processing service
processService := process.NewService(evaluator, "worker-01", operationRegistry)

// Process encrypted request
result, err := processService.ProcessRequest(request)

// Custom registry configuration
customRegistry := registry.NewService().
    Register(operation.OperationVariance, operation.NewVarianceOperation()).
    Register(operation.OperationMean, operation.NewMeanOperation())

processService := process.NewService(evaluator, "worker-02", customRegistry)
```

## Architecture Notes

- Depends only on domain types and port interfaces.
- Does not know implementation details of CKKS or Lattigo.
- Follows Dependency Inversion Principle through registry injection.
- Operations are managed externally via the registry service.
- Worker ID enables distributed system traceability.

## Supported Operations

Operations are managed by the registry service. The default registry includes:

- **Variance**: Statistical variance computation on encrypted data.

Additional operations can be registered through the registry service using `Register`.

## Related Components

- **Registry Service**: Manages available statistical operations.
- **Encrypt Service**: Client-side service that prepares encrypted requests.
- **Decrypt Service**: Client-side service that decrypts results.
- **Domain Entities**: Uses Request, Result, StatisticalOperation.
- **Ports**: HomomorphicEvaluator.
