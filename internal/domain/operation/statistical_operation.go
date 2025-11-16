package operation

// OperationType represents the type of statistical operation.
type OperationType string

const (
	OperationVariance          OperationType = "variance"
	OperationMean              OperationType = "mean"
	OperationStandardDeviation OperationType = "standard_deviation"
)

// StatisticalOperation defines the contract for statistical computations.
// Implementations must provide both plaintext and encrypted computation methods.
type StatisticalOperation interface {
	// Type returns the operation type identifier.
	Type() OperationType

	// ComputePlaintext performs the operation on unencrypted data.
	ComputePlaintext(data []float64, params OperationParams) (float64, error)

	// RequiredRotations returns the rotation keys needed for encrypted computation.
	RequiredRotations(dataSize int) []int

	// Validate checks if the operation can be performed with the given parameters.
	Validate(dataSize int, params OperationParams) error
}

// OperationParams contains parameters specific to each statistical operation.
type OperationParams struct {
	referenceValue float64
	additionalData map[string]interface{}
}

// NewOperationParams creates a new operation parameters object.
func NewOperationParams(referenceValue float64) OperationParams {
	return OperationParams{
		referenceValue: referenceValue,
		additionalData: make(map[string]interface{}),
	}
}

// ReferenceValue returns the reference value for the operation.
func (p OperationParams) ReferenceValue() float64 {
	return p.referenceValue
}

// SetAdditionalData stores additional operation-specific data.
func (p *OperationParams) SetAdditionalData(key string, value interface{}) {
	p.additionalData[key] = value
}

// GetAdditionalData retrieves additional operation-specific data.
func (p OperationParams) GetAdditionalData(key string) (interface{}, bool) {
	val, exists := p.additionalData[key]
	return val, exists
}
