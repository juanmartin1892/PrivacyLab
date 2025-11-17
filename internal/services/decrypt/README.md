# Decrypt Service# Decrypt Service



The decrypt service provides decryption and verification functionality for client applications. It processes encrypted computation results received from workers.The decrypt service provides decryption and verification functionality for client applications. It processes encrypted computation results received from workers.



## Purpose## Purpose



Transform encrypted computation results back into plaintext form and verify correctness against expected values.Transform encrypted computation results back into plaintext form and verify correctness against expected values.



## Responsibilities## Responsibilities



- Decrypt ciphertext results using the decryptor port.- Decrypt ciphertext results using the decryptor port.

- Decode encrypted results into plaintext values using the encoder port.- Decode encrypted results into plaintext values using the encoder port.

- Verify results against plaintext computation (for testing and validation).- Verify results against plaintext computation (for testing and validation).

- Compute plaintext reference values for verification purposes.

## Service API

## Service API

### Constructor

### Constructor

```go

NewService(```go

    decoder ports.Encoder,NewService(

    decryptor ports.Decryptor,    decoder ports.Encoder,

) *Service    decryptor ports.Decryptor,

```) *Service

```

Creates a new decryption service with the required infrastructure dependencies.

Creates a new decryption service with the required infrastructure dependencies.

### Methods

### Methods

- **DecryptResult**: Decrypts a computation result into plaintext value.

- **VerifyResult**: Verifies encrypted result matches expected plaintext value within tolerance.- **DecryptResult**: Decrypts a computation result into plaintext value.

- **VerifyResult**: Verifies encrypted result matches expected plaintext value within tolerance.

## Usage Example- **ComputePlaintext**: Computes plaintext reference value using operation's plaintext implementation.



```go## Usage Example

decryptService := decrypt.NewService(encoder, decryptor)

```go

result := computation.NewResult(encryptedData, opType, requestID, metadata)decryptService := decrypt.NewService(encoder, decryptor)

plaintext, err := decryptService.DecryptResult(result)

result := computation.NewResult(encryptedData, opType, requestID, metadata)

// Compute expected plaintext value using helper functionplaintext, err := decryptService.DecryptResult(result)

expected := computePlaintextVariance(referenceValue, population)

// Compute expected plaintext value

// Verify correctnessexpected, err := decryptService.ComputePlaintext(operation, data, params)

isValid, absoluteError, err := decryptService.VerifyResult(result, expected, tolerance)

```// Verify correctness

isValid, absoluteError, err := decryptService.VerifyResult(result, expected, tolerance)

## Architecture Notes```



- Depends only on domain types and port interfaces.## Architecture Notes

- Does not know implementation details of CKKS or Lattigo.

- Verification functionality useful for testing and debugging.- Depends only on domain types and port interfaces.

- Plaintext computation is delegated to example/test helper functions, not domain operations.- Does not know implementation details of CKKS or Lattigo.

- Verification functionality useful for testing and debugging.

## Related Components- Can compute plaintext reference values for any operation type.



- **Encrypt Service**: Counterpart service for encrypting data.## Related Components

- **Process Service**: Worker-side service that produces encrypted results.

- **Domain Entities**: Uses Result.- **Encrypt Service**: Counterpart service for encrypting data.

- **Ports**: Encoder (for decoding), Decryptor.- **Process Service**: Worker-side service that produces encrypted results.

- **Domain Entities**: Uses Result, Dataset, StatisticalOperation.
- **Ports**: Encoder (for decoding), Decryptor.
