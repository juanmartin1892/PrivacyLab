# Decrypt Service

The decrypt service provides decryption and verification functionality for client applications. It processes encrypted computation results received from workers.

## Purpose

Transform encrypted computation results back into plaintext form and verify correctness against expected values.

## Responsibilities

- Decrypt ciphertext results using the decryptor port.
- Decode encrypted results into plaintext values using the encoder port.
- Verify results against plaintext computation (for testing and validation).
- Compute plaintext reference values for verification purposes.

## Service API

### Constructor

```go
NewService(
    decoder ports.Encoder,
    decryptor ports.Decryptor,
) *Service
```

Creates a new decryption service with the required infrastructure dependencies.

### Methods

- **DecryptResult**: Decrypts a computation result into plaintext value.
- **VerifyResult**: Verifies encrypted result matches plaintext computation.
- **ComputePlaintext**: Computes plaintext reference value for verification.

## Usage Example

```go
decryptService := decrypt.NewService(encoder, decryptor)

result := computation.NewResult(encryptedData, opType, requestID, metadata)
plaintext, err := decryptService.DecryptResult(result)

// Verify correctness
isValid, err := decryptService.VerifyResult(result, datasets, singleValues, tolerance)
```

## Architecture Notes

- Depends only on domain types and port interfaces.
- Does not know implementation details of CKKS or Lattigo.
- Verification functionality useful for testing and debugging.
- Can compute plaintext reference values for any operation type.

## Related Components

- **Encrypt Service**: Counterpart service for encrypting data.
- **Process Service**: Worker-side service that produces encrypted results.
- **Domain Entities**: Uses Result, Dataset, StatisticalOperation.
- **Ports**: Encoder (for decoding), Decryptor.
