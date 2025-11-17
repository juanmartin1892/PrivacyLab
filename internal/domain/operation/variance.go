package operation

import (
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/errors"
)

const (
	// VarianceType is the standard identifier for variance operations.
	VarianceType OperationType = "variance"
)

// VarianceOperation implements the variance statistical operation.
// Variance formula: Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)
type VarianceOperation struct{}

// NewVarianceOperation creates a new variance operation.
func NewVarianceOperation() *VarianceOperation {
	return &VarianceOperation{}
}

// Type returns the operation type identifier.
func (v *VarianceOperation) Type() OperationType {
	return VarianceType
}

// ComputePlaintext calculates variance in plain text.
// This implementation calculates the variance of the data against a reference value.
// The reference value must be provided in the data slice as the first element,
// followed by the population data.
func (v *VarianceOperation) ComputePlaintext(data []float64, params OperationParams) (float64, error) {
	if len(data) < 2 {
		return 0, errors.ErrInsufficientData
	}

	x := data[0]
	population := data[1:]
	n := float64(len(population))
	sumSquaredDiff := 0.0

	for _, value := range population {
		diff := value - x
		sumSquaredDiff += diff * diff
	}

	return sumSquaredDiff / n, nil
}

// Validate checks if the variance operation can be performed with the given parameters.
func (v *VarianceOperation) Validate(dataSize int, params OperationParams) error {
	if dataSize == 0 {
		return fmt.Errorf("%w: cannot calculate variance on empty dataset", errors.ErrInsufficientData)
	}
	return nil
}
