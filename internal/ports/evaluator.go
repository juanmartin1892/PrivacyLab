package ports

import (
	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/operation"
)

// HomomorphicEvaluator defines the contract for performing homomorphic operations.
// This is a port that will be implemented by infrastructure adapters (e.g., CKKS adapter).
type HomomorphicEvaluator interface {
	// EvaluateOperation performs a statistical operation on encrypted data.
	// Returns encrypted result and metadata about the computation.
	EvaluateOperation(
		op operation.StatisticalOperation,
		request *computation.Request,
	) (*computation.Result, error)

	// Add performs homomorphic addition of two encrypted values.
	Add(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error)

	// Sub performs homomorphic subtraction.
	Sub(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error)

	// Mul performs homomorphic multiplication with relinearization.
	Mul(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error)

	// MulScalar performs multiplication by a plaintext scalar.
	MulScalar(ct *crypto.EncryptedData, scalar float64) (*crypto.EncryptedData, error)

	// Rotate performs slot rotation (requires rotation keys).
	Rotate(ct *crypto.EncryptedData, positions int) (*crypto.EncryptedData, error)

	// Rescale reduces the scale of a ciphertext after multiplication.
	Rescale(ct *crypto.EncryptedData) (*crypto.EncryptedData, error)

	// SumSlots sums n slots of a ciphertext into the first slot.
	SumSlots(ct *crypto.EncryptedData, n int) (*crypto.EncryptedData, error)
}
