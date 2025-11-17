package registry

import (
	"testing"

	"github.com/juanmartin/privacylab/internal/domain/operation"
)

func TestServiceRegister(t *testing.T) {
	service := NewService()
	mockOp := operation.NewVarianceOperation()

	service.Register(operation.VarianceType, mockOp)

	op, exists := service.Get(operation.VarianceType)
	if !exists {
		t.Fatal("expected operation to be registered")
	}
	if op == nil {
		t.Fatal("expected non-nil operation")
	}
}

func TestServiceGetNonExistent(t *testing.T) {
	service := NewService()

	_, exists := service.Get(operation.VarianceType)
	if exists {
		t.Fatal("expected operation to not exist")
	}
}

func TestServiceRegisterDefaults(t *testing.T) {
	service := NewService().RegisterDefaults()

	types := service.SupportedTypes()
	if len(types) == 0 {
		t.Fatal("expected at least one default operation")
	}

	_, exists := service.Get(operation.VarianceType)
	if !exists {
		t.Fatal("expected variance operation to be registered by default")
	}
}

func TestServiceChaining(t *testing.T) {
	service := NewService().
		Register(operation.VarianceType, operation.NewVarianceOperation())

	if len(service.SupportedTypes()) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(service.SupportedTypes()))
	}
}

func TestServiceSupportedTypes(t *testing.T) {
	service := NewService()

	types := service.SupportedTypes()
	if len(types) != 0 {
		t.Fatalf("expected 0 operations in empty registry, got %d", len(types))
	}

	service.Register(operation.VarianceType, operation.NewVarianceOperation())

	types = service.SupportedTypes()
	if len(types) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(types))
	}

	found := false
	for _, opType := range types {
		if opType == operation.VarianceType {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected variance operation in supported types")
	}
}

func TestServiceMultipleRegistrations(t *testing.T) {
	service := NewService()

	op1 := operation.NewVarianceOperation()
	service.Register(operation.VarianceType, op1)

	op2 := operation.NewVarianceOperation()
	service.Register(operation.VarianceType, op2)

	retrieved, exists := service.Get(operation.VarianceType)
	if !exists {
		t.Fatal("expected operation to exist")
	}

	if retrieved != op2 {
		t.Fatal("expected second registration to overwrite first")
	}
}
