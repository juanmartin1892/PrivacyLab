package operation

const (
	// VarianceType is the standard identifier for variance operations.
	VarianceType OperationType = "variance"
)

// varianceOperation is a minimal implementation for operation registration.
// Variance formula: Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)
// The actual computation logic is implemented in infrastructure adapters.
type varianceOperation struct{}

// NewVarianceOperation creates a new variance operation for registration.
func NewVarianceOperation() *varianceOperation {
	return &varianceOperation{}
}

// Type returns the operation type identifier.
func (v *varianceOperation) Type() OperationType {
	return VarianceType
}
