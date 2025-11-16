# Encrypt Service

The encrypt service provides encryption functionality for client applications. It coordinates between domain entities and cryptographic infrastructure to prepare encrypted computation requests.

## Purpose

Transform plaintext data and computation parameters into encrypted form, ready to be sent to a worker for homomorphic processing.

## Responsibilities

- Generate cryptographic key material through the key generator port.
- Encrypt datasets and individual values using the encryptor port.
- Encode plaintext data into CKKS format using the encoder port.
- Prepare computation requests with all necessary encrypted inputs and public keys.
- Support multiple operation types (variance, mean, etc.).

## Service API

### Constructor

```go
NewService(
    keyGen ports.KeyGenerator,
    encoder ports.Encoder,
    encryptor ports.Encryptor,
) *Service
```

Creates a new encryption service with the required infrastructure dependencies.

### Methods

- **PrepareRequest**: Prepares a generic encrypted computation request for any operation type.
- **EncryptDataset**: Encrypts a domain dataset into ciphertext.
- **EncryptValue**: Encrypts a single float64 value.
- **GetPublicKeySet**: Retrieves the public key material needed by workers.

## Usage Example

```go
encryptService := encrypt.NewService(keyGen, encoder, encryptor)

dataset := dataset.NewDataset([]float64{1.0, 2.0, 3.0})
request, err := encryptService.PrepareRequest(
    operation.OperationVariance,
    []*dataset.Dataset{dataset},
    nil,
    nil,
    "request-123",
)
```

## Architecture Notes

- Depends only on domain types and port interfaces.
- Does not know implementation details of CKKS or Lattigo.
- Supports the Open/Closed Principle by accepting any operation type.
- Thread-safe if underlying port implementations are thread-safe.

## Related Components

- **Decrypt Service**: Counterpart service for decrypting results.
- **Process Service**: Worker-side service that processes encrypted requests.
- **Domain Entities**: Uses Dataset, EncryptedData, CryptoParameters, Request.
- **Ports**: KeyGenerator, Encoder, Encryptor.
