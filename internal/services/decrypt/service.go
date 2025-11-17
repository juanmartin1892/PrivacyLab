package decrypt

import (
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/ports"
)

// Service handles decryption and verification of computation results.
// This service is used by the client to decrypt results received from the worker.
type Service struct {
	decoder   ports.Encoder
	decryptor ports.Decryptor
}

// NewService creates a new decryption service.
func NewService(decoder ports.Encoder, decryptor ports.Decryptor) *Service {
	return &Service{
		decoder:   decoder,
		decryptor: decryptor,
	}
}

// DecryptResult decrypts a computation result and extracts the value.
func (s *Service) DecryptResult(result *computation.Result) (float64, error) {
	// Decrypt the encrypted result
	value, err := s.decryptor.DecryptValue(result.EncryptedResult())
	if err != nil {
		return 0, fmt.Errorf("decryption failed: %w", err)
	}

	return value, nil
}

// VerifyResult verifies the encrypted result against plaintext computation.
func (s *Service) VerifyResult(
	result *computation.Result,
	plainTextResult float64,
	tolerance float64,
) (bool, float64, error) {
	// Decrypt encrypted result
	decryptedValue, err := s.DecryptResult(result)
	if err != nil {
		return false, 0, fmt.Errorf("failed to decrypt result: %w", err)
	}

	// Calculate absolute error
	absoluteError := abs(decryptedValue - plainTextResult)

	// Check if within tolerance
	isValid := absoluteError <= tolerance

	return isValid, absoluteError, nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
