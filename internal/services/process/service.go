package process

import (
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/domain/operation"
	"github.com/juanmartin/privacylab/internal/ports"
)

// Service handles encrypted computation requests.
// It coordinates between domain logic and cryptographic infrastructure.
type Service struct {
	evaluator  ports.HomomorphicEvaluator
	operations map[operation.OperationType]operation.StatisticalOperation
	workerID   string
}

// NewService creates a new encrypted processing service.
func NewService(evaluator ports.HomomorphicEvaluator, workerID string) *Service {
	return &Service{
		evaluator: evaluator,
		operations: map[operation.OperationType]operation.StatisticalOperation{
			operation.OperationVariance: operation.NewVarianceOperation(),
		},
		workerID: workerID,
	}
}

// ProcessRequest processes an encrypted computation request.
// This is the main entry point for the worker application.
func (s *Service) ProcessRequest(request *computation.Request) (*computation.Result, error) {
	// Get operation implementation
	op, exists := s.operations[request.OperationType()]
	if !exists {
		return nil, fmt.Errorf("unsupported operation: %s", request.OperationType())
	}

	// Validate operation can be performed
	dataSize := request.EncryptedInputs()[0].DataSize()
	params := operation.NewOperationParams(0) // Will be extracted from request parameters

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

// RegisterOperation allows adding new statistical operations dynamically.
func (s *Service) RegisterOperation(op operation.StatisticalOperation) {
	s.operations[op.Type()] = op
}

// SupportedOperations returns the list of currently supported operations.
func (s *Service) SupportedOperations() []operation.OperationType {
	types := make([]operation.OperationType, 0, len(s.operations))
	for t := range s.operations {
		types = append(types, t)
	}
	return types
}

// WorkerID returns the identifier of this worker.
func (s *Service) WorkerID() string {
	return s.workerID
}
