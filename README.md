# PrivacyLab

## Purpose

**PrivacyLab** is a hands-on lab designed to explore, prototype and demonstrate **advanced secure computing techniques**, combining:

- Homomorphic encryption
- Zero-knowledge proofs
- Isolation and separation of roles
- Secure key and identity architectures
- Multiparty computing in controlled contexts
- Verifiable end-to-end integrity
- Good engineering practices (SOLID + Hexagonal Architecture)

The goal is not to build a commercial product, but to create a controlled environment in which to systematically experiment and learn about these mechanisms.

## Vision

PrivacyLab's vision is to serve as a conceptual and practical basis for systems where:

- A **client** can execute operations on its data without disclosing it to the server.
- The **integrity** and **privacy** of data and results are demonstrably guaranteed.
- Results are **cryptographically verifiable**, without relying solely on trust in the infrastructure or the operator.
- The architecture is **modular, extensible and testable**, ready to integrate new secure computing techniques as they are explored.

## Current Scope

This project is in its **initial phase**. The current focus is on:

- Establishing the foundational architecture and project structure.
- Defining clear motivations and guiding principles.
- Setting up the development environment and tooling.
- Implementing initial prototypes and experiments with secure computing techniques.

**Status:** Active development with working examples.

### Implemented Features

#### Point-to-Point Variance Calculation with Homomorphic Encryption

A complete working example demonstrating privacy-preserving variance calculation using the CKKS homomorphic encryption scheme.

**Key Features:**
- **SIMD Batching**: All population values encrypted in a single ciphertext for enhanced privacy.
- **Hidden Population Size**: Prevents attackers from determining the number of individuals by counting ciphertexts.
- **Fully Homomorphic Operations**: Variance computed entirely on encrypted data without decryption.
- **Verifiable Results**: Plain-text comparison ensures computational correctness.

**Technical Highlights:**
- Uses Lattigo v6 library for homomorphic encryption.
- Implements CKKS scheme for approximate arithmetic on encrypted real numbers.
- Demonstrates rotation keys for SIMD operations.
- Achieves high precision with minimal error (relative error < 0.000001%).

**Location:** `examples/point-to-point/`

**Quick Start:**
```bash
cd examples/point-to-point
./run.sh
```

See `examples/point-to-point/README.md` for detailed documentation and execution instructions.

## Project Structure

```
privacyLab/
├── cmd/                    # Application entry points
│   ├── client/            # Client application (encryption, decryption)
│   └── worker/            # Worker application (homomorphic computation)
├── examples/              # Working examples and demonstrations
│   └── point-to-point/   # Variance calculation with homomorphic encryption
├── internal/              # Internal packages
│   ├── domain/           # Core business logic and domain models
│   │   ├── computation/  # Request/response model for client-worker communication
│   │   ├── crypto/       # Cryptographic configuration and key management
│   │   ├── dataset/      # Dataset entity representing numerical values
│   │   ├── errors/       # Domain-specific error definitions
│   │   └── operation/    # Statistical operations (variance, mean, etc.)
│   ├── ports/            # Interface definitions for infrastructure dependencies
│   ├── adapters/         # Infrastructure adapters implementing port interfaces
│   │   ├── lattigo/     # Lattigo v6 CKKS homomorphic encryption adapters
│   │   └── file/        # File-based dataset readers
│   └── services/         # Application services orchestrating domain operations
│       ├── encrypt/     # Client-side encryption and request preparation
│       ├── decrypt/     # Client-side decryption and verification
│       └── process/     # Worker-side computation processing
├── go.mod                # Go module dependencies
└── README.md             # This file
```

### Domain Layer

The domain layer (`internal/domain/`) contains the core business logic, completely independent of infrastructure:

- **Dataset** (`internal/domain/dataset/`): Immutable collections of validated numerical values with defensive copying and invariant enforcement
- **Operation** (`internal/domain/operation/`): Statistical operations (variance, etc.) with plaintext and encrypted computation support, including rotation key requirements
- **Crypto** (`internal/domain/crypto/`): Cryptographic parameters, keys (secret, public, relinearization, rotation), key sets, and encrypted data with metadata
- **Computation** (`internal/domain/computation/`): Request/response entities for distributed homomorphic computation with performance metadata tracking
- **Errors** (`internal/domain/errors/`): Centralized domain error definitions for validation, cryptographic configuration, and computation failures

All domain entities use only primitive types (no external library dependencies) and enforce invariants through validation.

**Architecture Principles:**
- Infrastructure independence: No dependencies on external libraries
- Immutability: Entities protect internal state through defensive copying
- Validation: All entities validate invariants at construction time
- SOLID principles: Clear separation of concerns and single responsibilities

**Client vs Worker Separation:**
- **Client (Trusted)**: Has plaintext data, owns secret keys, encrypts datasets, decrypts results, validates computations
- **Worker (Untrusted)**: Never sees plaintext, only has public/evaluation keys, performs homomorphic operations, returns encrypted results

See `internal/domain/README.md` for detailed domain architecture documentation.

### Adapter Layer

The adapter layer (`internal/adapters/`) implements port interfaces using external libraries and technologies:

#### Lattigo Adapters (`internal/adapters/lattigo/`)

Complete implementation of cryptographic operations using Lattigo v6 CKKS scheme:

- **KeyGeneratorAdapter**: Generates secret, public, relinearization, and rotation keys
- **EncoderAdapter**: Encodes plaintext values into CKKS plaintexts with SIMD batching
- **EncryptorAdapter**: Encrypts encoded plaintexts using public keys
- **DecryptorAdapter**: Decrypts ciphertexts and extracts values using secret keys
- **EvaluatorAdapter**: Performs homomorphic operations (add, sub, mul, rotate, rescale, sum slots)

**Features:**
- CKKS scheme for approximate arithmetic on encrypted real numbers
- SIMD batching for efficient multi-value encryption
- Support for rotation keys and slot operations
- Automatic level and scale management
- Serialization for network transmission and storage

See `internal/adapters/lattigo/README.md` for detailed documentation.

#### File Adapters (`internal/adapters/file/`)

Data reading adapters for loading datasets:

- **FileReaderAdapter**: Reads numerical datasets from text files and byte arrays
  - One value per line format
  - Comment support (lines starting with `#`)
  - Automatic validation and parsing

See `internal/adapters/file/README.md` for detailed documentation.

### Services Layer

The services layer (`internal/services/`) orchestrates domain logic and infrastructure components through application services:

#### Encrypt Service (`internal/services/encrypt/`)

Client-side service for data encryption and computation request preparation.

**Responsibilities:**
- Generate cryptographic key material
- Encrypt datasets and individual values
- Prepare computation requests with encrypted inputs
- Provide public key sets for workers

#### Decrypt Service (`internal/services/decrypt/`)

Client-side service for result decryption and verification.

**Responsibilities:**
- Decrypt computation results from workers
- Verify correctness against plaintext computation
- Compute reference values for validation
- Support result verification with configurable tolerance

#### Process Service (`internal/services/process/`)

Worker-side service for encrypted computation processing.

**Responsibilities:**
- Validate and execute computation requests
- Coordinate between operations and homomorphic evaluator
- Track performance metadata (time, operations, levels)
- Support dynamic operation registration
- Maintain worker identity for traceability

**Design Principles:**
- Single responsibility per service
- Dependency inversion through port interfaces
- Open/closed for extension through operation registration
- Comprehensive test coverage with mocks and integration tests

See `internal/services/README.md` for detailed service architecture and usage patterns.

#### Ports Layer

The ports layer (`internal/ports/`) defines interface contracts for infrastructure dependencies:

- **KeyGenerator**: Cryptographic key generation (secret, public, relinearization, rotation)
- **Encoder**: Plaintext encoding/decoding with SIMD batching support
- **Encryptor**: Public key encryption with metadata tracking
- **Decryptor**: Secret key decryption and value extraction
- **HomomorphicEvaluator**: Homomorphic operations on encrypted data (add, sub, mul, rotate, rescale, sum slots)
- **DataReader**: Dataset loading from external sources (files, bytes)

These interfaces ensure the domain layer remains independent of specific cryptographic libraries or data sources through dependency inversion.

See `internal/ports/README.md` for detailed interface specifications and usage patterns.

## Examples and Demonstrations

### 1. Point-to-Point Variance Calculation

Demonstrates privacy-preserving statistical computation using homomorphic encryption.

- **Location:** `examples/point-to-point/`
- **Technology:** CKKS homomorphic encryption (Lattigo v6)
- **Use Case:** Calculate variance of an input value against an encrypted population
- **Privacy Features:** SIMD batching, hidden population size, fully encrypted computation
- **Documentation:** See `examples/point-to-point/README.md`

**Run it:**
```bash
cd examples/point-to-point && ./run.sh
```

## Motivations

- Explore in a practical way how to combine homomorphic encryption, zero-knowledge proofs and MPC to **minimize information disclosure**
- Understand the **performance and complexity limits** of these techniques in a controlled environment
- Design and document **architectural patterns** to serve as a reference for future sensitive data computing systems
- Build a reusable knowledge base, with clear documentation, that facilitates the introduction of new use cases or new cryptographic primitives
- Apply **SOLID principles** and **Hexagonal Architecture** to create maintainable, testable, and extensible privacy-preserving systems

## Testing

The project includes comprehensive functional tests for all service layers to ensure code quality and consistency.

### Service Tests

All core services have complete test coverage with mocks and integration tests:

- **Encryption Service** (`internal/services/encrypt`): Tests for dataset/value encryption, request preparation, and public key management
- **Decryption Service** (`internal/services/decrypt`): Tests for result decryption, verification, and plaintext computation
- **Processing Service** (`internal/services/process`): Tests for encrypted computation, operation registration, and worker management

**Total Test Coverage:** 37 passing tests across all services

**Run all tests:**
```bash
go test ./internal/services/...
```

**Run tests for a specific service:**
```bash
go test ./internal/services/encrypt
go test ./internal/services/decrypt
go test ./internal/services/process
```

### Testing Strategy

The project employs comprehensive testing approaches across all layers:

- **Domain Layer**: Unit tests for validation logic, property-based tests for invariants, no mocks needed (pure domain logic)
- **Service Layer**: Functional tests with mock implementations of ports for isolation, integration tests for complete workflows
- **Adapter Layer**: Unit tests for individual functions, integration tests with external libraries

See individual package READMEs for detailed testing documentation.

## Extensibility

The architecture is designed to support future operations and features with minimal changes:

### Adding New Statistical Operations

1. Implement the `StatisticalOperation` interface in the domain layer
2. Define plaintext computation logic
3. Specify required rotation keys
4. Register the operation with the process service

No changes needed to existing domain entities, services, or adapters.

**Example operations to implement:**
- Mean (average) calculation
- Standard deviation
- Covariance between datasets
- Linear regression

### Adding New Encryption Schemes

1. Define new port interfaces if needed
2. Implement adapters for the new scheme (e.g., BFV for integers)
3. Create parameter configurations in the crypto domain
4. Inject new adapters into services

**Potential schemes:**
- BFV for integer arithmetic
- TFHE for fast bootstrapping
- BGV for modular arithmetic

### Adding New Infrastructure

1. Define port interface in `internal/ports/`
2. Implement adapter in `internal/adapters/`
3. Document usage in adapter README
4. Inject through service constructors

**Future adapters:**
- Network communication (gRPC, REST)
- Database storage
- Cloud key management (KMS, HSM)
- Distributed worker pools

## Author

**Author:** Juan Martín Pérez

Personal research and exploration project in secure computing.