# Services Layer

The services layer contains application services that orchestrate domain logic and infrastructure components. Each service has a specific responsibility in the homomorphic encryption workflow.

## Architecture

This layer follows the hexagonal architecture pattern:

- Services depend on domain entities and port interfaces.
- Services do NOT depend on concrete infrastructure implementations.
- Infrastructure adapters implement the ports and are injected via constructors.

## Service Packages

### Encrypt Service

**Location**: `internal/services/encrypt/`

**Responsibility**: Client-side encryption of data and computation request preparation.

**Main Operations**:
- Generate cryptographic keys.
- Encrypt datasets and values.
- Prepare computation requests.

**Dependencies**: KeyGenerator, Encoder, Encryptor ports.

### Decrypt Service

**Location**: `internal/services/decrypt/`

**Responsibility**: Client-side decryption of results and verification.

**Main Operations**:
- Decrypt computation results.
- Verify correctness against plaintext.
- Compute reference values.

**Dependencies**: Encoder, Decryptor ports.

### Process Service

**Location**: `internal/services/process/`

**Responsibility**: Worker-side processing of encrypted computation requests.

**Main Operations**:
- Validate operation requests.
- Execute homomorphic computations.
- Track performance metadata.

**Dependencies**: HomomorphicEvaluator port.

## Design Principles

### Single Responsibility

Each service has one clear purpose:
- Encrypt: Prepare encrypted requests.
- Decrypt: Process encrypted results.
- Process: Execute encrypted computations.

### Dependency Inversion

Services depend on abstractions (ports), not concrete implementations:

```go
type Service struct {
    evaluator ports.HomomorphicEvaluator  // Port, not *CKKSEvaluator
}
```

### Open/Closed

Services are open for extension (new operations) but closed for modification:

```go
processService.RegisterOperation(newOperation)
```

## Usage Flow

### Client Side

1. **Encrypt Service** prepares encrypted request.
2. Request sent to worker (via transport layer, not shown).
3. **Decrypt Service** processes encrypted result from worker.

### Worker Side

1. Receive encrypted request (via transport layer, not shown).
2. **Process Service** executes homomorphic computation.
3. Return encrypted result (via transport layer, not shown).

## Testing Strategy

All services include comprehensive functional tests with mocks and integration tests.

### Test Coverage

Each service has complete test coverage:

- **Encrypt Service Tests** (`encrypt/service_test.go`):
  - Service creation and initialization
  - Request preparation with datasets and single values
  - Dataset encryption with error handling
  - Value encryption with error handling
  - Public key set generation

- **Decrypt Service Tests** (`decrypt/service_test.go`):
  - Service creation and initialization
  - Result decryption with error handling
  - Result verification against plaintext
  - Plaintext computation
  - End-to-end integration tests

- **Process Service Tests** (`process/service_test.go`):
  - Service creation with default operations
  - Request processing with various scenarios
  - Dynamic operation registration
  - Supported operations listing
  - Worker identification
  - End-to-end integration tests

### Running Tests

Run all service tests:
```bash
go test ./internal/services/...
```

Run tests for a specific service:
```bash
go test ./internal/services/encrypt
go test ./internal/services/decrypt
go test ./internal/services/process
```

Run with verbose output:
```bash
go test -v ./internal/services/...
```

### Mock Implementations

Tests use mock implementations of ports for isolation:

```go
type mockEvaluator struct {
    evaluateOperationFunc func(op operation.StatisticalOperation, request *computation.Request) (*computation.Result, error)
}

func (m *mockEvaluator) EvaluateOperation(op operation.StatisticalOperation, request *computation.Request) (*computation.Result, error) {
    if m.evaluateOperationFunc != nil {
        return m.evaluateOperationFunc(op, request)
    }
    // Default behavior
}
```

This approach allows:
- Testing service logic in isolation.
- Simulating error conditions.
- Fast test execution without real cryptographic operations.
- Controlled and predictable test scenarios.

### Integration Tests

Each service includes integration tests that validate complete workflows:

```go
func TestServiceIntegration(t *testing.T) {
    service, encoder, decryptor := newTestService()
    // Test complete workflow
}
```

For detailed test documentation, see `docs/services-functional-tests-summary.md`.

## Extension Points

### Adding New Operations

1. Implement `operation.StatisticalOperation` interface in domain layer.
2. Register operation with process service: `processService.RegisterOperation(op)`.
3. No changes needed to service code.

### Adding New Ports

1. Define port interface in `internal/ports/`.
2. Update service constructor to receive new port.
3. Implement adapter in infrastructure layer.

## Related Documentation

- **Domain Layer**: `internal/domain/README.md`
- **Ports Layer**: `internal/ports/README.md`
- **Hexagonal Architecture**: `docs/hexagonal-architecture-implementation.md`
