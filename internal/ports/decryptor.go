package ports

import (
	"github.com/juanmartin/privacylab/internal/domain/crypto"
)

// Decryptor defines the contract for decrypting ciphertext data.
type Decryptor interface {
	// Decrypt decrypts an encrypted data object using the secret key.
	// Returns the encoded plaintext as bytes.
	Decrypt(encrypted *crypto.EncryptedData) ([]byte, error)

	// DecryptValue decrypts and decodes the first slot value.
	// Convenient for single-value results.
	DecryptValue(encrypted *crypto.EncryptedData) (float64, error)
}
