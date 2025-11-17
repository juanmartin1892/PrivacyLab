package registry

import (
	"github.com/juanmartin/privacylab/internal/domain/operation"
)

// Service manages the registration and retrieval of statistical operations.
// It provides a central registry for all available operations in the system.
type Service struct {
	operations map[operation.OperationType]operation.StatisticalOperation
}

// NewService creates an empty operation registry service.
func NewService() *Service {
	return &Service{
		operations: make(map[operation.OperationType]operation.StatisticalOperation),
	}
}

// Register adds a new operation to the registry.
// Returns the service to allow method chaining.
func (s *Service) Register(opType operation.OperationType, op operation.StatisticalOperation) *Service {
	s.operations[opType] = op
	return s
}

// Get retrieves an operation by type.
// Returns the operation and a boolean indicating if it exists.
func (s *Service) Get(opType operation.OperationType) (operation.StatisticalOperation, bool) {
	op, exists := s.operations[opType]
	return op, exists
}

// SupportedTypes returns all registered operation types.
func (s *Service) SupportedTypes() []operation.OperationType {
	types := make([]operation.OperationType, 0, len(s.operations))
	for t := range s.operations {
		types = append(types, t)
	}
	return types
}

// RegisterDefaults registers all standard operations.
// This provides a convenient way to get all built-in operations.
func (s *Service) RegisterDefaults() *Service {
	s.operations[operation.VarianceType] = operation.NewVarianceOperation()
	return s
}
