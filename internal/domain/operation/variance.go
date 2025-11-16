package operation

import (
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/errors"
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
	return OperationVariance
}

// ComputePlaintext calculates variance in plain text.
// The reference value in params is the value x to compare against the population.
func (v *VarianceOperation) ComputePlaintext(data []float64, params OperationParams) (float64, error) {
	if len(data) == 0 {
		return 0, errors.ErrInsufficientData
	}

	x := params.ReferenceValue()
	n := float64(len(data))
	sumSquaredDiff := 0.0

	for _, value := range data {
		diff := value - x
		sumSquaredDiff += diff * diff
	}

	return sumSquaredDiff / n, nil
}

// RequiredRotations returns the rotation keys needed for encrypted variance computation.
// For variance calculation with SIMD batching, we need rotations 1 to n-1 to sum all slots.
func (v *VarianceOperation) RequiredRotations(dataSize int) []int {
	rotations := make([]int, 0, dataSize-1)
	for i := 1; i < dataSize; i++ {
		rotations = append(rotations, i)
	}
	return rotations
}

// Validate checks if the variance operation can be performed with the given parameters.
func (v *VarianceOperation) Validate(dataSize int, params OperationParams) error {
	if dataSize == 0 {
		return fmt.Errorf("%w: cannot calculate variance on empty dataset", errors.ErrInsufficientData)
	}
	return nil
}
