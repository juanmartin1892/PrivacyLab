package process

import (
	"errors"
	"testing"
	"time"

	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/operation"
	"github.com/juanmartin/privacylab/internal/services/registry"
)

type mockHomomorphicEvaluator struct {
	evaluateOperationFunc func(op operation.StatisticalOperation, request *computation.Request) (*computation.Result, error)
	addFunc               func(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error)
	subFunc               func(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error)
	mulFunc               func(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error)
	mulScalarFunc         func(ct *crypto.EncryptedData, scalar float64) (*crypto.EncryptedData, error)
	rotateFunc            func(ct *crypto.EncryptedData, positions int) (*crypto.EncryptedData, error)
	rescaleFunc           func(ct *crypto.EncryptedData) (*crypto.EncryptedData, error)
	sumSlotsFunc          func(ct *crypto.EncryptedData, n int) (*crypto.EncryptedData, error)
}

func (m *mockHomomorphicEvaluator) EvaluateOperation(op operation.StatisticalOperation, request *computation.Request) (*computation.Result, error) {
	if m.evaluateOperationFunc != nil {
		return m.evaluateOperationFunc(op, request)
	}

	metadata := crypto.NewEncryptionMetadata(100, 2, 1<<40, 100, "test_hash")
	encryptedResult := crypto.NewEncryptedData([]byte("result_ciphertext"), metadata)

	resultMetadata := computation.NewResultMetadata(
		500*time.Millisecond,
		3,
		2,
		1,
		2,
		1<<40,
		"test-worker-1",
	)

	return computation.NewResult(encryptedResult, op.Type(), request.RequestID(), resultMetadata), nil
}

func (m *mockHomomorphicEvaluator) Add(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	if m.addFunc != nil {
		return m.addFunc(a, b)
	}
	metadata := crypto.NewEncryptionMetadata(100, 3, 1<<40, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("add_result"), metadata), nil
}

func (m *mockHomomorphicEvaluator) Sub(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	if m.subFunc != nil {
		return m.subFunc(a, b)
	}
	metadata := crypto.NewEncryptionMetadata(100, 3, 1<<40, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("sub_result"), metadata), nil
}

func (m *mockHomomorphicEvaluator) Mul(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	if m.mulFunc != nil {
		return m.mulFunc(a, b)
	}
	metadata := crypto.NewEncryptionMetadata(100, 2, 1<<40, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("mul_result"), metadata), nil
}

func (m *mockHomomorphicEvaluator) MulScalar(ct *crypto.EncryptedData, scalar float64) (*crypto.EncryptedData, error) {
	if m.mulScalarFunc != nil {
		return m.mulScalarFunc(ct, scalar)
	}
	metadata := crypto.NewEncryptionMetadata(100, 3, 1<<40, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("mulscalar_result"), metadata), nil
}

func (m *mockHomomorphicEvaluator) Rotate(ct *crypto.EncryptedData, positions int) (*crypto.EncryptedData, error) {
	if m.rotateFunc != nil {
		return m.rotateFunc(ct, positions)
	}
	metadata := crypto.NewEncryptionMetadata(100, 3, 1<<40, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("rotate_result"), metadata), nil
}

func (m *mockHomomorphicEvaluator) Rescale(ct *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	if m.rescaleFunc != nil {
		return m.rescaleFunc(ct)
	}
	metadata := crypto.NewEncryptionMetadata(100, 2, 1<<39, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("rescale_result"), metadata), nil
}

func (m *mockHomomorphicEvaluator) SumSlots(ct *crypto.EncryptedData, n int) (*crypto.EncryptedData, error) {
	if m.sumSlotsFunc != nil {
		return m.sumSlotsFunc(ct, n)
	}
	metadata := crypto.NewEncryptionMetadata(100, 3, 1<<40, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("sumslots_result"), metadata), nil
}

type mockStatisticalOperation struct {
	opType                operation.OperationType
	computePlaintextFunc  func(data []float64, params operation.OperationParams) (float64, error)
	requiredRotationsFunc func(dataSize int) []int
	validateFunc          func(dataSize int, params operation.OperationParams) error
}

func (m *mockStatisticalOperation) Type() operation.OperationType {
	if m.opType != "" {
		return m.opType
	}
	return operation.VarianceType
}

func (m *mockStatisticalOperation) ComputePlaintext(data []float64, params operation.OperationParams) (float64, error) {
	if m.computePlaintextFunc != nil {
		return m.computePlaintextFunc(data, params)
	}
	return 10.0, nil
}

func (m *mockStatisticalOperation) RequiredRotations(dataSize int) []int {
	if m.requiredRotationsFunc != nil {
		return m.requiredRotationsFunc(dataSize)
	}
	return []int{1, 2, 4}
}

func (m *mockStatisticalOperation) Validate(dataSize int, params operation.OperationParams) error {
	if m.validateFunc != nil {
		return m.validateFunc(dataSize, params)
	}
	return nil
}

func newTestService() (*Service, *mockHomomorphicEvaluator) {
	evaluator := &mockHomomorphicEvaluator{}
	reg := registry.NewService().RegisterDefaults()
	service := NewService(evaluator, "test-worker-1", reg)
	return service, evaluator
}

func createTestRequest(opType operation.OperationType, requestID string, dataSize int) *computation.Request {
	metadata := crypto.NewEncryptionMetadata(dataSize, 5, 1<<40, dataSize, "test_hash")
	encryptedData := crypto.NewEncryptedData([]byte("test_ciphertext"), metadata)

	request := computation.NewRequest(opType, []*crypto.EncryptedData{encryptedData}, requestID)
	request.SetParameter("data_size", dataSize)
	request.SetParameter("reference_value", 3.0)

	return request
}

func TestNewService(t *testing.T) {
	service, _ := newTestService()

	if service == nil {
		t.Fatal("expected service to be created")
	}

	if service.evaluator == nil {
		t.Error("expected evaluator to be set")
	}

	if service.registry == nil {
		t.Error("expected registry to be initialized")
	}

	if service.workerID != "test-worker-1" {
		t.Errorf("expected workerID to be 'test-worker-1', got %s", service.workerID)
	}

	supported := service.SupportedOperations()
	if len(supported) != 1 {
		t.Errorf("expected 1 default operation, got %d", len(supported))
	}

	foundVariance := false
	for _, opType := range supported {
		if opType == operation.VarianceType {
			foundVariance = true
			break
		}
	}
	if !foundVariance {
		t.Error("expected variance operation to be registered by default")
	}
}

func TestServiceProcessRequest(t *testing.T) {
	tests := []struct {
		name      string
		opType    operation.OperationType
		requestID string
		dataSize  int
		setupMock func(*mockHomomorphicEvaluator)
		wantError bool
		errorMsg  string
	}{
		{
			name:      "successful processing",
			opType:    operation.VarianceType,
			requestID: "req-001",
			dataSize:  100,
			setupMock: nil,
			wantError: false,
		},
		{
			name:      "unsupported operation",
			opType:    operation.OperationType("unsupported_op"),
			requestID: "req-002",
			dataSize:  100,
			setupMock: nil,
			wantError: true,
			errorMsg:  "unsupported operation",
		},
		{
			name:      "evaluation error",
			opType:    operation.VarianceType,
			requestID: "req-003",
			dataSize:  100,
			setupMock: func(m *mockHomomorphicEvaluator) {
				m.evaluateOperationFunc = func(op operation.StatisticalOperation, request *computation.Request) (*computation.Result, error) {
					return nil, errors.New("evaluation failed")
				}
			},
			wantError: true,
			errorMsg:  "evaluation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, evaluator := newTestService()

			if tt.setupMock != nil {
				tt.setupMock(evaluator)
			}

			request := createTestRequest(tt.opType, tt.requestID, tt.dataSize)

			result, err := service.ProcessRequest(request)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected result to be non-nil")
			}

			if result.OperationType() != tt.opType {
				t.Errorf("expected operation type %s, got %s", tt.opType, result.OperationType())
			}

			if result.RequestID() != tt.requestID {
				t.Errorf("expected request ID %s, got %s", tt.requestID, result.RequestID())
			}

			if result.Metadata().WorkerID() != "test-worker-1" {
				t.Errorf("expected worker ID 'test-worker-1', got %s", result.Metadata().WorkerID())
			}
		})
	}
}

func TestServiceSupportedOperations(t *testing.T) {
	service, _ := newTestService()

	supported := service.SupportedOperations()

	if len(supported) == 0 {
		t.Fatal("expected at least one supported operation")
	}

	foundVariance := false
	for _, opType := range supported {
		if opType == operation.VarianceType {
			foundVariance = true
			break
		}
	}

	if !foundVariance {
		t.Error("expected variance to be in supported operations")
	}
}

func TestServiceWorkerID(t *testing.T) {
	evaluator := &mockHomomorphicEvaluator{}
	reg := registry.NewService().RegisterDefaults()
	service := NewService(evaluator, "custom-worker-id", reg)

	if service.WorkerID() != "custom-worker-id" {
		t.Errorf("expected worker ID 'custom-worker-id', got %s", service.WorkerID())
	}
}

func TestServiceIntegration(t *testing.T) {
	service, evaluator := newTestService()

	callCount := 0
	evaluator.evaluateOperationFunc = func(op operation.StatisticalOperation, request *computation.Request) (*computation.Result, error) {
		callCount++

		if op.Type() != operation.VarianceType {
			t.Errorf("expected operation type %s, got %s", operation.VarianceType, op.Type())
		}

		if request.RequestID() != "integration-test" {
			t.Errorf("expected request ID 'integration-test', got %s", request.RequestID())
		}

		metadata := crypto.NewEncryptionMetadata(100, 2, 1<<40, 100, "test_hash")
		encryptedResult := crypto.NewEncryptedData([]byte("integrated_result"), metadata)

		resultMetadata := computation.NewResultMetadata(
			750*time.Millisecond,
			5,
			3,
			2,
			2,
			1<<40,
			"test-worker-1",
		)

		return computation.NewResult(encryptedResult, op.Type(), request.RequestID(), resultMetadata), nil
	}

	request := createTestRequest(operation.VarianceType, "integration-test", 100)

	result, err := service.ProcessRequest(request)
	if err != nil {
		t.Fatalf("ProcessRequest failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected evaluator to be called once, got %d calls", callCount)
	}

	if result == nil {
		t.Fatal("expected result to be non-nil")
	}

	if result.OperationType() != operation.VarianceType {
		t.Errorf("expected operation type %s, got %s", operation.VarianceType, result.OperationType())
	}

	if result.RequestID() != "integration-test" {
		t.Errorf("expected request ID 'integration-test', got %s", result.RequestID())
	}

	if result.Metadata().WorkerID() != "test-worker-1" {
		t.Errorf("expected worker ID 'test-worker-1', got %s", result.Metadata().WorkerID())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
