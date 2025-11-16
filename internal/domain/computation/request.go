package computation

import (
	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/operation"
)

// Request represents a computation request to be processed by the worker.
// It encapsulates all information needed to perform a calculation on encrypted data.
type Request struct {
	operationType   operation.OperationType
	encryptedInputs []*crypto.EncryptedData
	parameters      map[string]interface{}
	requestID       string
}

// NewRequest creates a new computation request.
func NewRequest(
	operationType operation.OperationType,
	encryptedInputs []*crypto.EncryptedData,
	requestID string,
) *Request {
	return &Request{
		operationType:   operationType,
		encryptedInputs: copyEncryptedDataSlice(encryptedInputs),
		parameters:      make(map[string]interface{}),
		requestID:       requestID,
	}
}

// OperationType returns the type of operation to perform.
func (r *Request) OperationType() operation.OperationType {
	return r.operationType
}

// EncryptedInputs returns a copy of the encrypted input data.
func (r *Request) EncryptedInputs() []*crypto.EncryptedData {
	return copyEncryptedDataSlice(r.encryptedInputs)
}

// RequestID returns the unique identifier for this request.
func (r *Request) RequestID() string {
	return r.requestID
}

// SetParameter sets an operation-specific parameter.
func (r *Request) SetParameter(key string, value interface{}) {
	r.parameters[key] = value
}

// GetParameter retrieves an operation-specific parameter.
func (r *Request) GetParameter(key string) (interface{}, bool) {
	val, exists := r.parameters[key]
	return val, exists
}

// Parameters returns all operation-specific parameters.
func (r *Request) Parameters() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range r.parameters {
		result[k] = v
	}
	return result
}

// DataSize returns the size of the dataset being processed.
// Extracts from parameters or from first encrypted input metadata.
func (r *Request) DataSize() int {
	if size, exists := r.parameters["data_size"]; exists {
		if intSize, ok := size.(int); ok {
			return intSize
		}
	}
	if len(r.encryptedInputs) > 0 {
		return r.encryptedInputs[0].DataSize()
	}
	return 0
}

// OperationParams converts request parameters to operation parameters.
func (r *Request) OperationParams() operation.OperationParams {
	refValue := 0.0
	if val, exists := r.parameters["reference_value"]; exists {
		if floatVal, ok := val.(float64); ok {
			refValue = floatVal
		}
	}
	return operation.NewOperationParams(refValue)
}

func copyEncryptedDataSlice(data []*crypto.EncryptedData) []*crypto.EncryptedData {
	result := make([]*crypto.EncryptedData, len(data))
	copy(result, data)
	return result
}
