package ports

import (
	"github.com/juanmartin/privacylab/internal/domain/crypto"
)

// KeyGenerator defines the contract for generating cryptographic keys.
type KeyGenerator interface {
	// GenerateKeys generates a complete key set for the given parameters.
	// Includes secret, public, relinearization, and rotation keys.
	GenerateKeys(params *crypto.Parameters, rotations []int) (*crypto.KeySet, error)

	// GenerateSecretKey generates only a secret key.
	GenerateSecretKey(params *crypto.Parameters) (*crypto.KeyMaterial, error)

	// GeneratePublicKey generates a public key from a secret key.
	GeneratePublicKey(params *crypto.Parameters, secretKey *crypto.KeyMaterial) (*crypto.KeyMaterial, error)

	// GenerateRelinearizationKey generates a relinearization key.
	GenerateRelinearizationKey(params *crypto.Parameters, secretKey *crypto.KeyMaterial) (*crypto.KeyMaterial, error)

	// GenerateRotationKeys generates rotation keys for specified positions.
	GenerateRotationKeys(params *crypto.Parameters, secretKey *crypto.KeyMaterial, positions []int) (map[int]*crypto.KeyMaterial, error)
}
