package decrypt

import (
	"errors"
	"math"
	"testing"

	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/operation"
)

type mockEncoder struct {
	encodeFunc       func(values []float64) ([]byte, error)
	encodeSingleFunc func(value float64, slots int) ([]byte, error)
	decodeFunc       func(encoded []byte, slots int) ([]float64, error)
	maxSlots         int
}

func (m *mockEncoder) Encode(values []float64) ([]byte, error) {
	if m.encodeFunc != nil {
		return m.encodeFunc(values)
	}
	return []byte("encoded"), nil
}

func (m *mockEncoder) EncodeSingle(value float64, slots int) ([]byte, error) {
	if m.encodeSingleFunc != nil {
		return m.encodeSingleFunc(value, slots)
	}
	return []byte("encoded_single"), nil
}

func (m *mockEncoder) Decode(encoded []byte, slots int) ([]float64, error) {
	if m.decodeFunc != nil {
		return m.decodeFunc(encoded, slots)
	}
	return []float64{1.0, 2.0, 3.0}, nil
}

func (m *mockEncoder) MaxSlots() int {
	if m.maxSlots > 0 {
		return m.maxSlots
	}
	return 8192
}

type mockDecryptor struct {
	decryptFunc      func(encrypted *crypto.EncryptedData) ([]byte, error)
	decryptValueFunc func(encrypted *crypto.EncryptedData) (float64, error)
}

func (m *mockDecryptor) Decrypt(encrypted *crypto.EncryptedData) ([]byte, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(encrypted)
	}
	return []byte("decrypted"), nil
}

func (m *mockDecryptor) DecryptValue(encrypted *crypto.EncryptedData) (float64, error) {
	if m.decryptValueFunc != nil {
		return m.decryptValueFunc(encrypted)
	}
	return 42.5, nil
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
	return operation.OperationVariance
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

func newTestService() (*Service, *mockEncoder, *mockDecryptor) {
	encoder := &mockEncoder{maxSlots: 8192}
	decryptor := &mockDecryptor{}
	service := NewService(encoder, decryptor)
	return service, encoder, decryptor
}

func createTestResult(value float64, opType operation.OperationType, requestID string) *computation.Result {
	metadata := crypto.NewEncryptionMetadata(100, 3, 1<<40, 100, "test_hash")
	encryptedData := crypto.NewEncryptedData([]byte("test_ciphertext"), metadata)

	resultMetadata := computation.NewResultMetadata(
		1000,
		5,
		3,
		2,
		3,
		1<<40,
		"worker-1",
	)

	return computation.NewResult(encryptedData, opType, requestID, resultMetadata)
}

func TestNewService(t *testing.T) {
	service, _, _ := newTestService()

	if service == nil {
		t.Fatal("expected service to be created")
	}

	if service.decoder == nil {
		t.Error("expected decoder to be set")
	}

	if service.decryptor == nil {
		t.Error("expected decryptor to be set")
	}
}

func TestServiceDecryptResult(t *testing.T) {
	tests := []struct {
		name          string
		expectedValue float64
		setupMock     func(*mockDecryptor)
		wantError     bool
		errorMsg      string
	}{
		{
			name:          "successful decryption",
			expectedValue: 42.5,
			setupMock: func(m *mockDecryptor) {
				m.decryptValueFunc = func(encrypted *crypto.EncryptedData) (float64, error) {
					return 42.5, nil
				}
			},
			wantError: false,
		},
		{
			name:          "decryption error",
			expectedValue: 0.0,
			setupMock: func(m *mockDecryptor) {
				m.decryptValueFunc = func(encrypted *crypto.EncryptedData) (float64, error) {
					return 0.0, errors.New("decryption failed")
				}
			},
			wantError: true,
			errorMsg:  "decryption failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, decryptor := newTestService()

			if tt.setupMock != nil {
				tt.setupMock(decryptor)
			}

			result := createTestResult(tt.expectedValue, operation.OperationVariance, "test-123")

			value, err := service.DecryptResult(result)

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

			if value != tt.expectedValue {
				t.Errorf("expected value %f, got %f", tt.expectedValue, value)
			}
		})
	}
}

func TestServiceVerifyResult(t *testing.T) {
	tests := []struct {
		name           string
		decryptedValue float64
		plaintextValue float64
		tolerance      float64
		expectedValid  bool
		expectedError  float64
		wantError      bool
		errorMsg       string
	}{
		{
			name:           "result within tolerance",
			decryptedValue: 10.0,
			plaintextValue: 10.05,
			tolerance:      0.1,
			expectedValid:  true,
			expectedError:  0.05,
			wantError:      false,
		},
		{
			name:           "result outside tolerance",
			decryptedValue: 10.0,
			plaintextValue: 11.0,
			tolerance:      0.5,
			expectedValid:  false,
			expectedError:  1.0,
			wantError:      false,
		},
		{
			name:           "exact match",
			decryptedValue: 42.5,
			plaintextValue: 42.5,
			tolerance:      0.001,
			expectedValid:  true,
			expectedError:  0.0,
			wantError:      false,
		},
		{
			name:           "decryption error",
			decryptedValue: 0.0,
			plaintextValue: 10.0,
			tolerance:      0.1,
			wantError:      true,
			errorMsg:       "failed to decrypt result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, decryptor := newTestService()

			decryptor.decryptValueFunc = func(encrypted *crypto.EncryptedData) (float64, error) {
				if tt.wantError {
					return 0.0, errors.New("decryption error")
				}
				return tt.decryptedValue, nil
			}

			result := createTestResult(tt.decryptedValue, operation.OperationVariance, "test-verify")

			isValid, absoluteError, err := service.VerifyResult(result, tt.plaintextValue, tt.tolerance)

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

			if isValid != tt.expectedValid {
				t.Errorf("expected isValid=%v, got %v", tt.expectedValid, isValid)
			}

			if math.Abs(absoluteError-tt.expectedError) > 1e-9 {
				t.Errorf("expected absolute error %f, got %f", tt.expectedError, absoluteError)
			}
		})
	}
}

func TestServiceComputePlaintext(t *testing.T) {
	tests := []struct {
		name          string
		data          []float64
		params        operation.OperationParams
		expectedValue float64
		wantError     bool
		errorMsg      string
	}{
		{
			name:          "successful computation",
			data:          []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			params:        operation.NewOperationParams(3.0),
			expectedValue: 2.0,
			wantError:     false,
		},
		{
			name:          "computation error",
			data:          []float64{},
			params:        operation.NewOperationParams(0.0),
			expectedValue: 0.0,
			wantError:     true,
			errorMsg:      "plaintext computation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _, _ := newTestService()

			mockOp := &mockStatisticalOperation{
				opType: operation.OperationVariance,
				computePlaintextFunc: func(data []float64, params operation.OperationParams) (float64, error) {
					if len(data) == 0 {
						return 0.0, errors.New("insufficient data")
					}
					return tt.expectedValue, nil
				},
			}

			result, err := service.ComputePlaintext(mockOp, tt.data, tt.params)

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

			if result != tt.expectedValue {
				t.Errorf("expected result %f, got %f", tt.expectedValue, result)
			}
		})
	}
}

func TestServiceIntegration(t *testing.T) {
	service, _, decryptor := newTestService()

	decryptedValue := 15.3
	plaintextValue := 15.25
	tolerance := 0.1

	decryptor.decryptValueFunc = func(encrypted *crypto.EncryptedData) (float64, error) {
		return decryptedValue, nil
	}

	result := createTestResult(decryptedValue, operation.OperationVariance, "integration-test")

	value, err := service.DecryptResult(result)
	if err != nil {
		t.Fatalf("DecryptResult failed: %v", err)
	}

	if value != decryptedValue {
		t.Errorf("expected decrypted value %f, got %f", decryptedValue, value)
	}

	isValid, absoluteError, err := service.VerifyResult(result, plaintextValue, tolerance)
	if err != nil {
		t.Fatalf("VerifyResult failed: %v", err)
	}

	if !isValid {
		t.Error("expected result to be valid")
	}

	expectedError := math.Abs(decryptedValue - plaintextValue)
	if math.Abs(absoluteError-expectedError) > 1e-9 {
		t.Errorf("expected error %f, got %f", expectedError, absoluteError)
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
