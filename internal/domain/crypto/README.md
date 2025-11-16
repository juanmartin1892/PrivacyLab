# Crypto Domain

This package defines the cryptographic configuration and key management entities used in homomorphic encryption operations.

## Purpose

Provides domain-level abstractions for:
- **Parameters**: Cryptographic scheme configuration (CKKS/BFV parameters)
- **KeyMaterial**: Individual cryptographic keys (secret, public, relinearization, rotation)
- **KeySet**: Collections of keys required for complete encryption workflows
- **EncryptedData**: Encrypted values with associated metadata

## Key Concepts

### Parameters
Encapsulates the configuration of a homomorphic encryption scheme:
- Polynomial degree (LogN)
- Modulus chains (LogQ, LogP)
- Scaling factors (LogDefaultScale)
- Security level (128, 192, or 256 bits)

All parameters are validated at construction time to ensure cryptographic security.

### KeyMaterial
Represents individual cryptographic keys as raw bytes, maintaining independence from specific cryptographic libraries. Each key includes:
- Type identification (secret, public, relinearization, rotation)
- Raw key data (as bytes)
- Configuration metadata (parameter compatibility, rotation positions)

### KeySet
Manages collections of related keys needed for client or worker operations. Supports:
- Public/secret key pairs (client only has secret key)
- Relinearization keys (for multiplication depth management)
- Rotation keys (for SIMD slot operations)

### EncryptedData
Represents ciphertext with metadata required for correct homomorphic operations:
- Raw ciphertext bytes
- Current multiplication level
- Scaling factor
- SIMD slots used
- Parameter compatibility hash

## Design Decisions

- All cryptographic data stored as primitive types (bytes, integers)
- No dependencies on external cryptographic libraries in domain layer
- Defensive copying to ensure immutability
- Parameter validation enforces security requirements
- Compatibility checking prevents mismatched operations
