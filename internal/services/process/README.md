# Process Service

The process service handles encrypted computation requests on the worker side. It coordinates between domain operation logic and homomorphic evaluation infrastructure.

## Purpose

Execute statistical operations on encrypted data without decrypting it, using homomorphic encryption capabilities.

## Responsibilities

- Validate that requested operations can be performed on provided data.
- Coordinate between operation definitions and homomorphic evaluator.
- Track computation performance (time, operations, memory).
- Support dynamic registration of new statistical operations.
- Maintain worker identity for result traceability.

## Service API

### Constructor

```go
NewService(
    evaluator ports.HomomorphicEvaluator,
    workerID string,
) *Service
```

Creates a new processing service with the required infrastructure dependencies.

### Methods

- **ProcessRequest**: Processes an encrypted computation request and returns encrypted result.
- **RegisterOperation**: Dynamically adds support for new statistical operations.
- **SupportedOperations**: Returns list of currently supported operation types.
- **WorkerID**: Returns the identifier of this worker instance.

## Usage Example

```go
processService := process.NewService(evaluator, "worker-01")

result, err := processService.ProcessRequest(request)

// Add support for new operation
processService.RegisterOperation(operation.NewMeanOperation())
```

## Architecture Notes

- Depends only on domain types and port interfaces.
- Does not know implementation details of CKKS or Lattigo.
- Supports Open/Closed Principle through operation registration.
- Can be extended with new operations at runtime.
- Worker ID enables distributed system traceability.

## Supported Operations

By default, the service supports:

- **Variance**: Statistical variance computation on encrypted data.

Additional operations can be registered dynamically using `RegisterOperation`.

## Related Components

- **Encrypt Service**: Client-side service that prepares encrypted requests.
- **Decrypt Service**: Client-side service that decrypts results.
- **Domain Entities**: Uses Request, Result, StatisticalOperation.
- **Ports**: HomomorphicEvaluator.
