package encrypt

import (
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/dataset"
	"github.com/juanmartin/privacylab/internal/domain/operation"
	"github.com/juanmartin/privacylab/internal/ports"
)

// Service handles the encryption of data and preparation of computation requests.
// This service is used by the client to encrypt sensitive data before sending to the worker.
type Service struct {
	encoder   ports.Encoder
	encryptor ports.Encryptor
	keySet    *crypto.KeySet
	params    *crypto.Parameters
}

// NewService creates a new encryption service.
func NewService(
	encoder ports.Encoder,
	encryptor ports.Encryptor,
	keySet *crypto.KeySet,
	params *crypto.Parameters,
) *Service {
	return &Service{
		encoder:   encoder,
		encryptor: encryptor,
		keySet:    keySet,
		params:    params,
	}
}

// PrepareRequest prepares an encrypted computation request.
// The client specifies the operation type, datasets to encrypt, single values to encrypt,
// and any additional parameters needed.
func (s *Service) PrepareRequest(
	operationType operation.OperationType,
	datasets []*dataset.Dataset,
	singleValues []float64,
	parameters map[string]interface{},
	requestID string,
) (*computation.Request, error) {
	var encryptedInputs []*crypto.EncryptedData

	// Encrypt all datasets
	for i, data := range datasets {
		encrypted, err := s.encryptDataset(data)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt dataset %d: %w", i, err)
		}
		encryptedInputs = append(encryptedInputs, encrypted)
	}

	// Encrypt all single values (replicated across slots)
	for i, value := range singleValues {
		encrypted, err := s.encryptSingleValue(value)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt single value %d: %w", i, err)
		}
		encryptedInputs = append(encryptedInputs, encrypted)
	}

	// Create computation request
	request := computation.NewRequest(
		operationType,
		encryptedInputs,
		requestID,
	)

	// Add all provided parameters
	for key, value := range parameters {
		request.SetParameter(key, value)
	}

	return request, nil
}

// EncryptDataset encrypts a dataset and returns the encrypted data.
// Useful when the client wants to manage encryption separately.
func (s *Service) EncryptDataset(data *dataset.Dataset) (*crypto.EncryptedData, error) {
	return s.encryptDataset(data)
}

// EncryptValue encrypts a single value and returns the encrypted data.
// Useful when the client wants to manage encryption separately.
func (s *Service) EncryptValue(value float64) (*crypto.EncryptedData, error) {
	return s.encryptSingleValue(value)
}

// GetPublicKeySet returns a key set suitable for sending to the worker.
// It excludes the secret key for security.
func (s *Service) GetPublicKeySet() *crypto.KeySet {
	publicKeySet := crypto.NewKeySet("")
	publicKeySet.SetPublicKey(s.keySet.PublicKey())
	publicKeySet.SetRelinearizationKey(s.keySet.RelinearizationKey())

	// Add all rotation keys
	for _, pos := range s.keySet.RotationPositions() {
		if key, exists := s.keySet.RotationKey(pos); exists {
			publicKeySet.AddRotationKey(pos, key)
		}
	}

	return publicKeySet
}

// encryptDataset encrypts a complete dataset using SIMD batching.
func (s *Service) encryptDataset(data *dataset.Dataset) (*crypto.EncryptedData, error) {
	// Encode dataset values
	encoded, err := s.encoder.Encode(data.Values())
	if err != nil {
		return nil, fmt.Errorf("encoding failed: %w", err)
	}

	// Encrypt encoded data
	encrypted, err := s.encryptor.Encrypt(encoded)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	return encrypted, nil
}

// encryptSingleValue encrypts a single value, replicating it across all slots.
func (s *Service) encryptSingleValue(value float64) (*crypto.EncryptedData, error) {
	// Encode single value (replicated across slots)
	encoded, err := s.encoder.EncodeSingle(value, s.encoder.MaxSlots())
	if err != nil {
		return nil, fmt.Errorf("encoding failed: %w", err)
	}

	// Encrypt encoded value
	encrypted, err := s.encryptor.Encrypt(encoded)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	return encrypted, nil
}
