package dataset

import (
	"fmt"
	"math"

	"github.com/juanmartin/privacylab/internal/domain/errors"
)

// Dataset represents a collection of numerical values to be processed.
// It encapsulates the raw data and provides validation.
type Dataset struct {
	values []float64
	size   int
}

// NewDataset creates a new dataset from a slice of values.
// Returns an error if the values are invalid or empty.
func NewDataset(values []float64) (*Dataset, error) {
	if len(values) == 0 {
		return nil, errors.ErrEmptyDataset
	}

	if err := validateValues(values); err != nil {
		return nil, err
	}

	return &Dataset{
		values: copyValues(values),
		size:   len(values),
	}, nil
}

// Values returns a copy of the dataset values to prevent external modifications.
func (d *Dataset) Values() []float64 {
	return copyValues(d.values)
}

// Size returns the number of values in the dataset.
func (d *Dataset) Size() int {
	return d.size
}

// ValueAt returns the value at the specified index.
// Returns an error if the index is out of bounds.
func (d *Dataset) ValueAt(index int) (float64, error) {
	if index < 0 || index >= d.size {
		return 0, fmt.Errorf("index %d out of bounds [0, %d)", index, d.size)
	}
	return d.values[index], nil
}

// validateValues checks if all values in the slice are valid numbers.
func validateValues(values []float64) error {
	for i, v := range values {
		if isInvalidNumber(v) {
			return fmt.Errorf("%w: invalid value at index %d: %v", errors.ErrInvalidDataValue, i, v)
		}
	}
	return nil
}

// isInvalidNumber checks if a float64 is NaN or infinite.
func isInvalidNumber(v float64) bool {
	return math.IsNaN(v) || math.IsInf(v, 0)
}

// copyValues creates a defensive copy of a slice.
func copyValues(values []float64) []float64 {
	result := make([]float64, len(values))
	copy(result, values)
	return result
}
