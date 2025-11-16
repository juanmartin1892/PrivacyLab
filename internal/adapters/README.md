# Adapters

This directory contains infrastructure adapters that implement the port interfaces defined in `internal/ports`.

## Architecture

The adapter layer follows the Hexagonal Architecture pattern, where adapters implement port interfaces to connect the domain layer with external technologies and libraries.

```
Domain Layer (internal/domain)
      ↓
Port Interfaces (internal/ports)
      ↓
Adapters (internal/adapters)
      ↓
External Libraries (Lattigo, file I/O, etc.)
```

## Available Adapters

### Lattigo Adapters (`lattigo/`)

Implements cryptographic ports using the Lattigo v6 library for homomorphic encryption.

**Implemented ports:**
- `ports.KeyGenerator` - Cryptographic key generation
- `ports.Encoder` - Plaintext encoding/decoding
- `ports.Encryptor` - Public key encryption
- `ports.Decryptor` - Secret key decryption
- `ports.HomomorphicEvaluator` - Homomorphic operations

**Scheme:** CKKS (Cheon-Kim-Kim-Song)

See [lattigo/README.md](lattigo/README.md) for detailed documentation.

### File Adapters (`file/`)

Implements data reading ports for loading datasets from files.

**Implemented ports:**
- `ports.DataReader` - Read datasets from files and byte arrays

**Supported formats:**
- Plain text (one value per line)
- Comments and empty lines

See [file/README.md](file/README.md) for detailed documentation.

## Creating New Adapters

To create a new adapter:

1. **Identify the port interface** you want to implement from `internal/ports`
2. **Create a new package** under `internal/adapters/` (e.g., `adapters/mytech/`)
3. **Implement the interface** using your chosen technology
4. **Add constructor functions** that return the port interface type
5. **Document your adapter** with a README.md file

Example structure:

```
internal/adapters/
├── mytech/
│   ├── README.md
│   ├── encoder.go          # Implements ports.Encoder
│   ├── encryptor.go        # Implements ports.Encryptor
│   └── key_generator.go    # Implements ports.KeyGenerator
```

### Adapter Guidelines

**DO:**
- Return port interface types from constructors
- Validate inputs and return descriptive errors
- Document public API and usage examples
- Follow domain error conventions
- Use defensive copying for mutable data
- Maintain infrastructure independence in the domain

**DON'T:**
- Expose external library types in public API
- Mix business logic with adapter code
- Bypass port interfaces
- Store mutable state without synchronization
- Hard-code configuration values

## Dependency Management

Adapters are the only layer allowed to depend on external libraries. The domain layer must never import infrastructure packages.

```
✓ adapters/lattigo → github.com/tuneinsight/lattigo/v6
✓ adapters/lattigo → internal/ports
✓ adapters/lattigo → internal/domain

✗ domain → github.com/tuneinsight/lattigo/v6
✗ domain → adapters/lattigo
✗ ports → adapters/lattigo
```

## Testing

Each adapter should include:
- Unit tests for individual functions
- Integration tests with the external library
- Example usage in documentation

Run adapter tests:
```bash
go test ./internal/adapters/...
```

## Future Adapters

Potential adapters to implement:

### Cryptographic Adapters
- **BFV adapter** - Integer homomorphic encryption scheme
- **TFHE adapter** - Fast bootstrapping for boolean circuits
- **Seal adapter** - Microsoft SEAL library integration

### Data Adapters
- **Database adapter** - Read datasets from SQL databases
- **Cloud storage adapter** - Read from S3, GCS, Azure Blob
- **Streaming adapter** - Real-time data ingestion

### Network Adapters
- **gRPC adapter** - Remote procedure calls
- **REST adapter** - HTTP-based communication
- **WebSocket adapter** - Bidirectional streaming

### Storage Adapters
- **Key storage** - Secure key management (HSM, KMS)
- **Result storage** - Persist computation results
- **Cache adapter** - In-memory caching layer

## Configuration

Adapters should accept configuration through:
1. Constructor parameters (preferred)
2. Configuration structs
3. Environment variables (for deployment)

Avoid global state and singletons.

## Performance

Consider these optimization strategies:
- Object pooling for frequently allocated objects
- Lazy initialization for expensive resources
- Caching for immutable data
- Parallel processing where safe
- Profiling to identify bottlenecks

## Security

Adapters handling cryptographic material must:
- Never log sensitive data (keys, plaintexts)
- Use secure random number generation
- Implement proper key lifecycle management
- Follow cryptographic best practices
- Validate all inputs

## Maintenance

When updating adapters:
- Maintain backward compatibility when possible
- Version breaking changes clearly
- Update documentation and examples
- Run full test suite
- Update dependency versions carefully
