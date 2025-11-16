package errors

import "errors"

var (
	// Dataset errors
	ErrEmptyDataset     = errors.New("dataset cannot be empty")
	ErrInvalidDataValue = errors.New("dataset contains invalid values")

	// Crypto errors
	ErrInvalidCryptoParams = errors.New("invalid cryptographic parameters")
	ErrIncompatibleKeys    = errors.New("keys are incompatible with parameters")
	ErrMissingKey          = errors.New("required key is missing")

	// Operation errors
	ErrInvalidOperation     = errors.New("invalid operation parameters")
	ErrUnsupportedOperation = errors.New("operation not supported")
	ErrInsufficientData     = errors.New("insufficient data for operation")

	// Computation errors
	ErrEncryptionFailed = errors.New("encryption failed")
	ErrDecryptionFailed = errors.New("decryption failed")
	ErrEvaluationFailed = errors.New("homomorphic evaluation failed")
)
