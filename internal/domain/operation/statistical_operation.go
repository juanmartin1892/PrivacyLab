package operation

// OperationType represents the type of statistical operation.
// This is a string-based identifier that allows dynamic registration
// of operations without modifying the domain layer.
type OperationType string

// StatisticalOperation defines the contract for statistical computations.
// This interface focuses purely on identifying the operation type.
// Implementation details (validation, plaintext computation, encryption logic)
// are delegated to infrastructure adapters.
type StatisticalOperation interface {
	// Type returns the operation type identifier.
	Type() OperationType
}

// OperationParams contains parameters specific to each statistical operation.
type OperationParams struct {
	data map[string]interface{}
}

// NewOperationParams creates a new operation parameters object.
func NewOperationParams() OperationParams {
	return OperationParams{
		data: make(map[string]interface{}),
	}
}

// Set stores an operation-specific parameter.
func (p *OperationParams) Set(key string, value interface{}) {
	p.data[key] = value
}

// Get retrieves an operation-specific parameter.
func (p OperationParams) Get(key string) (interface{}, bool) {
	val, exists := p.data[key]
	return val, exists
}
