package lattigo

import "github.com/juanmartin/privacylab/internal/domain/operation"

// RotationMapper maps statistical operations to their required rotation keys.
// This is an infrastructure concern specific to homomorphic encryption with SIMD batching.
type RotationMapper struct {
	mappings map[operation.OperationType]RotationStrategy
}

// RotationStrategy defines how to compute required rotations for an operation.
type RotationStrategy func(dataSize int) []int

// NewRotationMapper creates a new rotation mapper with default strategies.
func NewRotationMapper() *RotationMapper {
	mapper := &RotationMapper{
		mappings: make(map[operation.OperationType]RotationStrategy),
	}

	// Register default rotation strategies
	mapper.Register(operation.VarianceType, VarianceRotationStrategy)

	return mapper
}

// Register adds a rotation strategy for a specific operation type.
func (rm *RotationMapper) Register(opType operation.OperationType, strategy RotationStrategy) {
	rm.mappings[opType] = strategy
}

// GetRequiredRotations returns the rotation keys needed for an operation.
// Returns nil if no specific rotations are required or the operation is not registered.
func (rm *RotationMapper) GetRequiredRotations(opType operation.OperationType, dataSize int) []int {
	strategy, exists := rm.mappings[opType]
	if !exists {
		return nil
	}
	return strategy(dataSize)
}

// VarianceRotationStrategy computes rotations needed for variance calculation.
// For variance with SIMD batching, we need rotations 1 to n-1 to sum all slots.
func VarianceRotationStrategy(dataSize int) []int {
	if dataSize <= 1 {
		return []int{}
	}
	rotations := make([]int, 0, dataSize-1)
	for i := 1; i < dataSize; i++ {
		rotations = append(rotations, i)
	}
	return rotations
}
