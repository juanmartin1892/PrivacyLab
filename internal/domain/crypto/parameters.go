package crypto

import (
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/errors"
)

// SchemeType represents the type of homomorphic encryption scheme.
type SchemeType string

const (
	SchemeCKKS SchemeType = "CKKS"
	SchemeBFV  SchemeType = "BFV"
)

// Parameters represents the cryptographic configuration for homomorphic encryption.
// All fields use primitive types to maintain domain independence.
type Parameters struct {
	scheme          SchemeType
	logN            int
	logQ            []int
	logP            []int
	logDefaultScale int
	securityLevel   int
}

// NewParameters creates a new cryptographic parameters configuration.
func NewParameters(
	scheme SchemeType,
	logN int,
	logQ []int,
	logP []int,
	logDefaultScale int,
	securityLevel int,
) (*Parameters, error) {
	if err := validateParameters(scheme, logN, logQ, logP, logDefaultScale, securityLevel); err != nil {
		return nil, err
	}

	return &Parameters{
		scheme:          scheme,
		logN:            logN,
		logQ:            copyIntSlice(logQ),
		logP:            copyIntSlice(logP),
		logDefaultScale: logDefaultScale,
		securityLevel:   securityLevel,
	}, nil
}

// Scheme returns the encryption scheme type.
func (p *Parameters) Scheme() SchemeType {
	return p.scheme
}

// LogN returns the log of the polynomial degree.
func (p *Parameters) LogN() int {
	return p.logN
}

// LogQ returns a copy of the ciphertext modulus chain.
func (p *Parameters) LogQ() []int {
	return copyIntSlice(p.logQ)
}

// LogP returns a copy of the auxiliary modulus chain.
func (p *Parameters) LogP() []int {
	return copyIntSlice(p.logP)
}

// LogDefaultScale returns the default scaling factor.
func (p *Parameters) LogDefaultScale() int {
	return p.logDefaultScale
}

// SecurityLevel returns the security level in bits.
func (p *Parameters) SecurityLevel() int {
	return p.securityLevel
}

// MaxSlots returns the maximum number of SIMD slots available.
func (p *Parameters) MaxSlots() int {
	return 1 << (p.logN - 1)
}

// MaxLevel returns the maximum multiplication depth.
func (p *Parameters) MaxLevel() int {
	return len(p.logQ) - 1
}

func validateParameters(
	scheme SchemeType,
	logN int,
	logQ []int,
	logP []int,
	logDefaultScale int,
	securityLevel int,
) error {
	if scheme != SchemeCKKS && scheme != SchemeBFV {
		return fmt.Errorf("unsupported scheme: %s", scheme)
	}

	if logN < 10 || logN > 17 {
		return fmt.Errorf("%w: logN must be between 10 and 17, got %d", errors.ErrInvalidCryptoParams, logN)
	}

	if len(logQ) == 0 {
		return fmt.Errorf("%w: logQ cannot be empty", errors.ErrInvalidCryptoParams)
	}

	if logDefaultScale <= 0 {
		return fmt.Errorf("%w: logDefaultScale must be positive, got %d", errors.ErrInvalidCryptoParams, logDefaultScale)
	}

	if securityLevel != 128 && securityLevel != 192 && securityLevel != 256 {
		return fmt.Errorf("%w: security level must be 128, 192, or 256 bits, got %d", errors.ErrInvalidCryptoParams, securityLevel)
	}

	return nil
}

func copyIntSlice(slice []int) []int {
	result := make([]int, len(slice))
	copy(result, slice)
	return result
}
