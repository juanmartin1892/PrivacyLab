# Lattigo Adapters

This directory contains adapter implementations for homomorphic encryption operations using the Lattigo v6 library.

## Overview

The adapters in this package implement the ports defined in `internal/ports` using the CKKS (Cheon-Kim-Kim-Song) homomorphic encryption scheme from Lattigo v6. These adapters bridge the domain layer with the cryptographic infrastructure.

## Components

### RotationMapper

Maps statistical operations to their required rotation keys for homomorphic encryption with SIMD batching.

**Features:**
- Separates cryptographic infrastructure from domain logic
- Provides extensible strategies for different operations
- Supports registration of custom rotation strategies
- Includes built-in strategies for standard operations

**Usage:**
```go
// Create mapper with default strategies
mapper := lattigo.NewRotationMapper()

// Get rotations for variance operation
rotations := mapper.GetRequiredRotations(operation.VarianceType, 100)

// Use rotations for key generation
keySet, _ := keyGen.GenerateKeys(params, rotations)

// Register custom strategy
customStrategy := func(dataSize int) []int {
    // Return required rotation positions
    return []int{1, 2, 4, 8}
}
mapper.Register("custom_mean", customStrategy)
```

### KeyGeneratorAdapter

Implements `ports.KeyGenerator` for generating cryptographic keys.

**Features:**
- Generates secret keys for encryption/decryption
- Generates public keys for encryption
- Generates relinearization keys for efficient multiplication
- Generates rotation keys for SIMD slot operations

**Usage:**
```go
params, _ := crypto.NewParameters(
    crypto.SchemeCKKS,
    14,                           // logN
    []int{55, 40, 40, 40, 40, 40}, // logQ
    []int{61},                     // logP
    40,                            // logDefaultScale
    128,                           // security level
)

keygen, _ := lattigo.NewKeyGeneratorAdapter(params)

rotations := []int{1, 2, 4, 8, 16, 32, 64}
keySet, _ := keygen.GenerateKeys(params, rotations)
```

### EncoderAdapter

Implements `ports.Encoder` for encoding plaintext values into CKKS plaintexts.

**Features:**
- Batch encoding of multiple values into SIMD slots
- Single value encoding with replication across slots
- Decoding of encoded plaintexts back to float64 values

**Usage:**
```go
encoder, _ := lattigo.NewEncoderAdapter(params)

values := []float64{1.5, 2.7, 3.9, 4.2}
encoded, _ := encoder.Encode(values)

singleValue := 5.0
encodedSingle, _ := encoder.EncodeSingle(singleValue, encoder.MaxSlots())
```

### EncryptorAdapter

Implements `ports.Encryptor` for encrypting encoded plaintexts.

**Features:**
- Public key encryption
- Metadata tracking (level, scale, slots used)
- Encrypted data serialization

**Usage:**
```go
encryptor, _ := lattigo.NewEncryptorAdapter(params, keySet.PublicKey())

encoded, _ := encoder.Encode(values)
encrypted, _ := encryptor.Encrypt(encoded)
```

### DecryptorAdapter

Implements `ports.Decryptor` for decrypting ciphertexts.

**Features:**
- Secret key decryption
- Single value extraction from first slot
- Full slot decoding

**Usage:**
```go
decryptor, _ := lattigo.NewDecryptorAdapter(params, keySet.SecretKey())

decoded, _ := decryptor.Decrypt(encrypted)
singleValue, _ := decryptor.DecryptValue(encrypted)
```

### EvaluatorAdapter

Implements `ports.HomomorphicEvaluator` for homomorphic operations.

**Features:**
- Statistical operation evaluation (variance, mean, etc.)
- Homomorphic addition and subtraction
- Homomorphic multiplication with relinearization
- Scalar multiplication
- Slot rotation
- Rescaling after multiplication
- Slot summation for aggregation

**Usage:**
```go
evaluator, _ := lattigo.NewEvaluatorAdapter(params, keySet, "worker-1")

result1, _ := evaluator.Add(encryptedA, encryptedB)
result2, _ := evaluator.Mul(encryptedA, encryptedB)
result3, _ := evaluator.MulScalar(encryptedA, 2.5)
rotated, _ := evaluator.Rotate(encryptedA, 1)
summed, _ := evaluator.SumSlots(encryptedA, dataSize)
```

## CKKS Parameters

Typical parameter configurations:

**Small dataset (< 1000 values):**
```go
params, _ := crypto.NewParameters(
    crypto.SchemeCKKS,
    14,                           // 16384 slots
    []int{55, 40, 40, 40, 40, 40},
    []int{61},
    40,
    128,
)
```

**Medium dataset (1000-10000 values):**
```go
params, _ := crypto.NewParameters(
    crypto.SchemeCKKS,
    15,                           // 32768 slots
    []int{55, 45, 45, 45, 45, 45},
    []int{61},
    45,
    128,
)
```

**Large dataset (> 10000 values):**
```go
params, _ := crypto.NewParameters(
    crypto.SchemeCKKS,
    16,                           // 65536 slots
    []int{60, 50, 50, 50, 50, 50},
    []int{61},
    50,
    128,
)
```

## Rotation Keys and SIMD Operations

The `RotationMapper` component manages rotation key requirements for different statistical operations. This separates domain logic from cryptographic infrastructure.

### Using RotationMapper

For variance calculation and other statistical operations that require slot rotations:

```go
// Create rotation mapper
mapper := lattigo.NewRotationMapper()

// Get required rotations for specific operation
rotations := mapper.GetRequiredRotations(operation.VarianceType, dataSize)

// Generate keys with those rotations
keySet, _ := keygen.GenerateKeys(params, rotations)
```

### Built-in Rotation Strategies

**Variance:** Requires rotations 1 to n-1 for summing all slots
```go
// For dataset of size 5
rotations := mapper.GetRequiredRotations(operation.VarianceType, 5)
// Returns: [1, 2, 3, 4]
```

### Custom Rotation Strategies

Register custom strategies for new operations:

```go
// Define custom operation rotation requirements
meanStrategy := func(dataSize int) []int {
    // Power-of-2 rotations for efficient tree-based summation
    rotations := []int{}
    for i := 1; i < dataSize; i *= 2 {
        rotations = append(rotations, i)
    }
    return rotations
}

// Register the strategy
mapper.Register(operation.MeanType, meanStrategy)
```

### Alternative: Power-of-2 Rotations

For large datasets, power-of-2 rotations are more efficient:

```go
// More efficient for large datasets using tree-based summation
rotations := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}
```

## Serialization

All adapters use `encoding/gob` for serializing Lattigo objects to byte arrays. This allows:
- Storing keys and ciphertexts to disk
- Transmitting encrypted data over network
- Maintaining domain independence

## Error Handling

Adapters follow Go best practices:
- Return descriptive errors with context
- Wrap underlying Lattigo errors with `fmt.Errorf` and `%w`
- Validate inputs before processing

## Performance Considerations

1. **Key generation:** Expensive operation, generate keys once and reuse
2. **Rotation keys:** Only generate the rotations you need
3. **Rescaling:** Always rescale after multiplication to maintain precision
4. **Slot usage:** Pack as many values as possible into SIMD slots
5. **Level management:** Monitor multiplication depth to avoid running out of levels

## Security

- Uses recommended CKKS parameters from Lattigo
- Supports 128, 192, and 256-bit security levels
- Keys are never exposed outside the adapter layer
- All cryptographic operations are delegated to Lattigo

## Dependencies

- `github.com/tuneinsight/lattigo/v6/core/rlwe` - Core RLWE primitives
- `github.com/tuneinsight/lattigo/v6/schemes/ckks` - CKKS scheme implementation

## Testing

To test the adapters, see the examples in `examples/point-to-point` which demonstrate:
- Complete encryption workflow
- Homomorphic variance calculation
- Key generation and management
