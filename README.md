# PrivacyLab

**PrivacyLab** is a hands-on lab for exploring advanced secure computing techniques through practical experimentation with homomorphic encryption, zero-knowledge proofs, and multiparty computation.

## Vision

Enable systems where:
- Clients execute operations on data without disclosing it to servers
- Data integrity and privacy are cryptographically guaranteed
- Results are verifiable without trusting infrastructure
- Architecture remains modular, extensible, and testable

## Current Status

**Active development** with working examples. Initial phase focuses on foundational architecture and prototypes.

### Working Example: Variance Calculation

Privacy-preserving variance calculation using CKKS homomorphic encryption:
- SIMD batching: all values in single ciphertext
- Hidden population size prevents inference attacks
- Fully homomorphic computation without decryption
- High precision (relative error < 0.000001%)

**Quick Start:**
```bash
cd examples/point-to-point
./run.sh
```

See `examples/point-to-point/README.md` for details.

## Architecture

Hexagonal architecture with clear separation of concerns:

```
privacyLab/
├── cmd/              # Application entry points (client, worker)
├── examples/         # Working demonstrations
├── internal/
│   ├── domain/      # Core business logic (dataset, operations, crypto, computation)
│   ├── ports/       # Interface definitions for infrastructure
│   ├── adapters/    # Infrastructure implementations (Lattigo CKKS, file I/O)
│   └── services/    # Application orchestration (encrypt, decrypt, process)
```

### Key Layers

**Domain** (`internal/domain/`): Pure business logic with no infrastructure dependencies
- Dataset, Operation, Crypto primitives, Computation models, Error definitions
- Immutable entities with invariant validation
- Client (trusted) vs Worker (untrusted) separation

**Ports** (`internal/ports/`): Interface contracts for infrastructure
- KeyGenerator, Encoder, Encryptor, Decryptor, HomomorphicEvaluator, DataReader

**Adapters** (`internal/adapters/`): Infrastructure implementations
- Lattigo v6 CKKS: Complete homomorphic encryption implementation
- File readers: Dataset loading with validation

**Services** (`internal/services/`): Application orchestration
- Encrypt: Client-side key generation and data encryption
- Decrypt: Result decryption and verification
- Process: Worker-side encrypted computation

See package READMEs for detailed documentation.

## Testing

**37 passing tests** across all service layers with comprehensive coverage.

**Run all tests:**
```bash
go test ./internal/services/...
```

**Strategy:**
- Domain: Unit tests for validation, property-based tests for invariants
- Services: Functional tests with mocks, integration tests for workflows
- Adapters: Unit and integration tests with external libraries

## Extensibility

### Add Statistical Operations
1. Implement `StatisticalOperation` interface
2. Register with process service

### Add Encryption Schemes
1. Define port interfaces
2. Implement adapters (BFV, TFHE, BGV)
3. Inject into services

### Add Infrastructure
1. Define port in `internal/ports/`
2. Implement adapter in `internal/adapters/`
3. Inject through constructors

## Motivations

- Practical exploration of privacy-preserving techniques
- Understanding performance and complexity limits
- Documenting architectural patterns for sensitive data systems
- Building reusable knowledge base with clear documentation
- Applying SOLID principles and Hexagonal Architecture

## Author

Juan Martín Pérez - Personal research project in secure computing.