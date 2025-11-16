# Ports Package

This package defines the interfaces (ports) that connect the domain layer to the infrastructure layer, following the Hexagonal Architecture pattern.

## Purpose

Ports establish clear contracts between the domain and external systems:
- **Dependency Inversion**: Domain depends on abstractions (ports), not concrete implementations
- **Testability**: Easy to mock ports for testing domain logic
- **Flexibility**: Swap implementations without changing domain code
- **Clear Boundaries**: Explicit separation between business logic and infrastructure

## Port Interfaces

### HomomorphicEvaluator

Performs operations on encrypted data without decryption.

**File**: `evaluator.go`

**Responsibilities**:
- Execute statistical operations on encrypted data
- Basic homomorphic arithmetic (add, sub, mul, mul scalar)
- SIMD operations (rotate, sum slots)
- Scale management (rescale after multiplication)

**Implementation**: `internal/adapters/ckks/evaluator_adapter.go`

**Usage Context**: Worker application performing computations on encrypted data

---

### Encoder

Converts between plaintext values and encoded representations suitable for encryption.

**File**: `encoder.go`

**Responsibilities**:
- Encode plaintext values into polynomial representations
- Decode encrypted results back to plaintext
- Manage SIMD batching and slot packing
- Support single-value replication across slots

**Implementation**: `internal/adapters/ckks/encoder_adapter.go`

**Usage Context**: Both client (encoding before encryption) and worker (if needed for plaintext operations)

---

### Encryptor

Encrypts plaintext data using public key cryptography.

**File**: `encryptor.go`

**Responsibilities**:
- Encrypt encoded plaintext to ciphertext
- Manage encryption randomness
- Attach metadata to encrypted data
- Support batch encryption

**Implementation**: `internal/adapters/ckks/encryptor_adapter.go`

**Usage Context**: Client application encrypting sensitive data before sending to worker

---

### Decryptor

Decrypts ciphertext back to plaintext using secret key.

**File**: `decryptor.go`

**Responsibilities**:
- Decrypt ciphertext to encoded plaintext
- Extract single values from decrypted results
- Validate decryption correctness
- Handle decryption errors gracefully

**Implementation**: `internal/adapters/ckks/decryptor_adapter.go`

**Usage Context**: Client application decrypting results received from worker

---

### KeyGenerator

Generates all cryptographic keys needed for homomorphic encryption.

**File**: `key_generator.go`

**Responsibilities**:
- Generate secret and public key pairs
- Create relinearization keys for multiplication depth
- Generate rotation keys for SIMD operations
- Assemble complete key sets

**Implementation**: `internal/adapters/ckks/key_generator_adapter.go`

**Usage Context**: Client application during initialization

---

### DataReader

Reads datasets from external sources.

**File**: `data_reader.go`

**Responsibilities**:
- Load data from files
- Parse numeric values
- Validate input format
- Convert to domain Dataset objects

**Implementation**: `internal/adapters/file/dataset_reader.go`

**Usage Context**: Client application loading input data

## Design Principles

### Interface Segregation
Each port defines a focused set of related operations. Clients depend only on the interfaces they need.

### Technology Agnostic
Ports define capabilities without prescribing implementation details. The domain doesn't know about Lattigo, CKKS, or file formats.

### Explicit Dependencies
Application services declare their dependencies through port interfaces, making requirements clear.

### Testability
Mock implementations of ports enable isolated testing of domain and application logic.

## Dependency Flow

```
┌─────────────────┐
│   Application   │  (orchestrates use cases)
│    Services     │
└────────┬────────┘
         │ depends on
         ▼
┌─────────────────┐
│      Ports      │  (interfaces/contracts)
│   (this pkg)    │
└────────┬────────┘
         │ implemented by
         ▼
┌─────────────────┐
│    Adapters     │  (concrete implementations)
│  (CKKS, File)   │
└────────┬────────┘
         │ uses
         ▼
┌─────────────────┐
│    External     │  (Lattigo, OS file system)
│   Libraries     │
└─────────────────┘
```

The domain never directly depends on external libraries. All infrastructure access flows through ports.

## Usage Example

```go
// Application service declares dependencies via ports
type EncryptDataService struct {
    encoder   ports.Encoder
    encryptor ports.Encryptor
    reader    ports.DataReader
}

// Concrete adapters injected at runtime
service := &EncryptDataService{
    encoder:   ckks.NewEncoderAdapter(params),
    encryptor: ckks.NewEncryptorAdapter(params, keys),
    reader:    file.NewDatasetReader(),
}

// Domain logic remains clean and testable
data, _ := service.reader.ReadFromFile("data.txt")
encoded, _ := service.encoder.Encode(data.Values())
encrypted, _ := service.encryptor.Encrypt(encoded)
```

## Adding New Ports

When adding infrastructure capabilities:

1. Define the port interface in this package
2. Keep methods focused and minimal
3. Use domain types in signatures (not infrastructure types)
4. Document the purpose and usage context
5. Implement in appropriate adapter package
6. Update this README with the new port

## Implementation Packages

- **CKKS Adapters**: `internal/adapters/ckks/` - Lattigo-based implementations
- **File Adapters**: `internal/adapters/file/` - File system operations
- **Future**: Network adapters, database adapters, etc.
