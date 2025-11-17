# Registry Service

The registry service provides centralized management of statistical operations in the system. It acts as a registry pattern implementation for operation discovery and retrieval.

## Purpose

Manage the registration and retrieval of statistical operations, allowing the system to be extended with new operations without modifying existing code.

## Responsibilities

- Register statistical operations with their corresponding operation types.
- Retrieve operations by type for execution.
- Track all supported operation types in the system.
- Provide default operation configurations.

## Service API

### Constructor

```go
NewService() *Service
```

Creates an empty operation registry service. Operations must be registered before use.

### Methods

- **Register**: Adds a new operation to the registry. Returns self for method chaining.
- **Get**: Retrieves an operation by type. Returns operation and existence flag.
- **SupportedTypes**: Returns all registered operation types.
- **RegisterDefaults**: Registers all standard built-in operations.

## Usage Example

```go
// Create registry with default operations
registry := registry.NewService().RegisterDefaults()

// Create registry with custom operations
registry := registry.NewService().
    Register(operation.OperationVariance, operation.NewVarianceOperation()).
    Register(operation.OperationMean, operation.NewMeanOperation())

// Retrieve an operation
op, exists := registry.Get(operation.OperationVariance)
if !exists {
    return fmt.Errorf("unsupported operation")
}

// List all supported operations
types := registry.SupportedTypes()
```

## Architecture Notes

- Follows the Registry pattern for managing operation instances.
- Supports method chaining for fluent configuration.
- Enables Open/Closed Principle: extend without modification.
- Facilitates Dependency Injection for operation management.
- No dependencies on infrastructure or external systems.

## Design Benefits

### SOLID Principles

- **Single Responsibility**: Only manages operation registration and retrieval.
- **Open/Closed**: New operations can be added without modifying the registry.
- **Dependency Inversion**: Operations are injected, not created internally.

### Flexibility

- Easy to create different registry configurations for different contexts.
- Simplifies testing by allowing injection of mock operations.
- Supports both default and custom operation sets.

## Integration

The registry service is typically used by the process service to determine which operations are available for execution:

```go
registry := registry.NewService().RegisterDefaults()
processService := process.NewService(evaluator, workerID, registry)
```

## Related Components

- **Process Service**: Main consumer of the registry service.
- **Statistical Operations**: Domain entities managed by the registry.
- **Operation Types**: Domain value objects used as registry keys.
