package ports

import (
	"github.com/juanmartin/privacylab/internal/domain/crypto"
)

// Encryptor defines the contract for encrypting plaintext data.
type Encryptor interface {
	// Encrypt encrypts encoded plaintext data using the public key.
	// Returns encrypted data with associated metadata.
	Encrypt(encoded []byte) (*crypto.EncryptedData, error)

	// EncryptWithMetadata encrypts data and attaches custom metadata.
	EncryptWithMetadata(encoded []byte, metadata crypto.EncryptionMetadata) (*crypto.EncryptedData, error)
}
