package crypto

// EncryptedData represents encrypted data ready for homomorphic operations.
// It uses raw bytes to remain independent of any specific cryptographic library.
type EncryptedData struct {
	ciphertext []byte
	metadata   EncryptionMetadata
}

// EncryptionMetadata contains information about how the data was encrypted.
type EncryptionMetadata struct {
	dataSize      int
	level         int
	scale         float64
	slotsUsed     int
	parameterHash string
}

// NewEncryptedData creates a new encrypted data object.
func NewEncryptedData(ciphertext []byte, metadata EncryptionMetadata) *EncryptedData {
	return &EncryptedData{
		ciphertext: copyBytes(ciphertext),
		metadata:   metadata,
	}
}

// Ciphertext returns a copy of the encrypted data bytes.
func (e *EncryptedData) Ciphertext() []byte {
	return copyBytes(e.ciphertext)
}

// Metadata returns the encryption metadata.
func (e *EncryptedData) Metadata() EncryptionMetadata {
	return e.metadata
}

// DataSize returns the original size of the encrypted data.
func (e *EncryptedData) DataSize() int {
	return e.metadata.dataSize
}

// Level returns the current multiplication level of the ciphertext.
func (e *EncryptedData) Level() int {
	return e.metadata.level
}

// Scale returns the current scaling factor.
func (e *EncryptedData) Scale() float64 {
	return e.metadata.scale
}

// SlotsUsed returns the number of SIMD slots used.
func (e *EncryptedData) SlotsUsed() int {
	return e.metadata.slotsUsed
}

// IsCompatibleWith checks if two encrypted data objects can be operated together.
func (e *EncryptedData) IsCompatibleWith(other *EncryptedData) bool {
	return e.metadata.parameterHash == other.metadata.parameterHash
}

// NewEncryptionMetadata creates a new encryption metadata object.
func NewEncryptionMetadata(dataSize, level int, scale float64, slotsUsed int, parameterHash string) EncryptionMetadata {
	return EncryptionMetadata{
		dataSize:      dataSize,
		level:         level,
		scale:         scale,
		slotsUsed:     slotsUsed,
		parameterHash: parameterHash,
	}
}

func copyBytes(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}
