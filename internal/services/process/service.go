package process

import (
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/domain/operation"
	"github.com/juanmartin/privacylab/internal/ports"
	"github.com/juanmartin/privacylab/internal/services/registry"
)

// Service handles encrypted computation requests.
// It coordinates between domain logic and cryptographic infrastructure.
type Service struct {
	evaluator ports.HomomorphicEvaluator
	registry  *registry.Service
	workerID  string
}

// NewService creates a new encrypted processing service.
// The registry parameter defines which operations are supported.
// Use registry.NewService().RegisterDefaults() for standard operations.
func NewService(evaluator ports.HomomorphicEvaluator, workerID string, operationRegistry *registry.Service) *Service {
	if operationRegistry == nil {
		operationRegistry = registry.NewService()
	}

	return &Service{
		evaluator: evaluator,
		registry:  operationRegistry,
		workerID:  workerID,
	}
}

// ProcessRequest processes an encrypted computation request.
// This is the main entry point for the worker application.
func (s *Service) ProcessRequest(request *computation.Request) (*computation.Result, error) {
	// Get operation implementation
	op, exists := s.registry.Get(request.OperationType())
	if !exists {
		return nil, fmt.Errorf("unsupported operation: %s", request.OperationType())
	}

	// Validate all encrypted inputs have consistent data sizes
	inputs := request.EncryptedInputs()
	if len(inputs) == 0 {
		return nil, fmt.Errorf("request has no encrypted inputs")
	}

	dataSize := inputs[0].DataSize()
	for i, input := range inputs {
		if input.DataSize() != dataSize {
			return nil, fmt.Errorf("input %d has inconsistent data size: expected %d, got %d", i, dataSize, input.DataSize())
		}
	}

	// Create operation parameters from request
	params := operation.NewOperationParams()

	// Copy parameters from request to operation params
	for key, value := range request.Parameters() {
		params.Set(key, value)
	}

	if err := op.Validate(dataSize, params); err != nil {
		return nil, fmt.Errorf("operation validation failed: %w", err)
	}

	// Delegate to evaluator (infrastructure) to perform homomorphic computation
	// The evaluator returns a complete result with accurate metadata including
	// multiplications and rotations used
	result, err := s.evaluator.EvaluateOperation(op, request)
	if err != nil {
		return nil, fmt.Errorf("evaluation failed: %w", err)
	}

	return result, nil
}

// SupportedOperations returns the list of currently supported operations.
func (s *Service) SupportedOperations() []operation.OperationType {
	return s.registry.SupportedTypes()
}

// WorkerID returns the identifier of this worker.
func (s *Service) WorkerID() string {
	return s.workerID
}
