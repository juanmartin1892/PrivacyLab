# Domain Errors

This package defines all error types used throughout the PrivacyLab domain layer.

## Purpose

Centralize error definitions to ensure consistent error handling across all domain entities and operations. By defining errors at the domain level, we maintain clear separation between business logic failures and infrastructure-specific errors.

## Error Categories

- **Dataset errors**: Validation failures related to input data
- **Crypto errors**: Issues with cryptographic parameter configuration and key management
- **Operation errors**: Problems during statistical operation execution
- **Computation errors**: Failures during encryption, decryption, or homomorphic evaluation

## Usage

Import this package when implementing domain entities or ports that need to return domain-specific errors.

```go
import "github.com/juanmartin/privacyLab/internal/domain/errors"

if len(values) == 0 {
    return nil, errors.ErrEmptyDataset
}
```
