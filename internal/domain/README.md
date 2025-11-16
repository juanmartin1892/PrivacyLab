# Domain Layer

This layer contains the core business logic and domain entities of the PrivacyLab system, completely independent of infrastructure concerns.

## Architecture Principles

The domain layer follows these principles:
- **Infrastructure Independence**: No dependencies on external libraries or frameworks
- **Immutability**: Domain entities protect their internal state through defensive copying
- **Validation**: All entities validate their invariants at construction time
- **Primitive Types**: Uses only Go primitive types (bytes, integers, floats, strings)
- **SOLID Principles**: Clear separation of concerns and single responsibilities

## Package Organization

### `/dataset`
Defines the Dataset entity representing collections of numerical values to be processed.
- Validates data integrity (no NaN or Inf values)
- Provides immutable access to values
- Ensures non-empty datasets

### `/operation`
Defines statistical operations that can be performed on datasets.
- **StatisticalOperation**: Interface for all operations
- **VarianceOperation**: Variance calculation implementation
- Supports both plaintext and encrypted computation
- Specifies required cryptographic resources

### `/crypto`
Defines cryptographic configuration and key management entities.
- **Parameters**: HE scheme configuration (CKKS/BFV)
- **KeyMaterial**: Individual cryptographic keys
- **KeySet**: Collections of keys for client/worker
- **EncryptedData**: Ciphertext with metadata

### `/computation`
Defines the request-response model for client-worker communication.
- **Request**: Computation requests with encrypted inputs
- **Result**: Computation results with performance metadata
- Supports distributed computation workflows

### `/errors`
Centralized error definitions for the domain layer.
- Dataset validation errors
- Cryptographic configuration errors
- Operation execution errors
- Computation errors

## Domain Model Overview

```
┌─────────────┐
│   Dataset   │ ──────────┐
└─────────────┘           │
                          │
                          ▼
              ┌───────────────────────┐
              │ StatisticalOperation  │
              └───────────────────────┘
                          │
                          │ uses
                          ▼
              ┌───────────────────────┐
              │    EncryptedData      │
              └───────────────────────┘
                          │
                          │ created by
                          ▼
              ┌───────────────────────┐
              │      Parameters       │
              │      KeyMaterial      │
              └───────────────────────┘
                          │
                          │ used in
                          ▼
              ┌───────────────────────┐
              │  Request / Result     │
              └───────────────────────┘
```

## Client vs Worker Separation

### Client (Trusted Component)
- Has access to plaintext data
- Owns secret keys
- Creates encryption parameters
- Generates all cryptographic keys
- Encrypts datasets
- Decrypts results
- Validates computations

### Worker (Untrusted Component)
- Never sees plaintext data
- Only has public and evaluation keys
- Performs homomorphic operations
- Returns encrypted results
- Tracks computational cost

## Extensibility

The domain is designed to support future operations and features:
- New statistical operations (mean, standard deviation, covariance)
- Additional encryption schemes (BFV for integer arithmetic)
- Multi-party computation protocols
- Distributed worker pools

To add a new statistical operation:
1. Implement the `StatisticalOperation` interface
2. Define plaintext computation logic
3. Specify required rotation keys
4. No changes needed to existing domain entities

## Testing Strategy

Domain entities should be tested with:
- Unit tests for validation logic
- Property-based tests for invariants
- Round-trip tests (serialize/deserialize if applicable)
- Compatibility tests (encryption parameter combinations)
- No mocks needed (pure domain logic)
