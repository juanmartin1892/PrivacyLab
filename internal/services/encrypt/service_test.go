package encrypt

import (
	"errors"
	"testing"

	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/dataset"
	"github.com/juanmartin/privacylab/internal/domain/operation"
)

// Mock implementations for testing

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

type mockEncryptor struct {
	encryptFunc             func(encoded []byte) (*crypto.EncryptedData, error)
	encryptWithMetadataFunc func(encoded []byte, metadata crypto.EncryptionMetadata) (*crypto.EncryptedData, error)
}

func (m *mockEncryptor) Encrypt(encoded []byte) (*crypto.EncryptedData, error) {
	if m.encryptFunc != nil {
		return m.encryptFunc(encoded)
	}
	metadata := crypto.NewEncryptionMetadata(len(encoded), 5, 1<<40, 100, "test_hash")
	return crypto.NewEncryptedData([]byte("ciphertext"), metadata), nil
}

func (m *mockEncryptor) EncryptWithMetadata(encoded []byte, metadata crypto.EncryptionMetadata) (*crypto.EncryptedData, error) {
	if m.encryptWithMetadataFunc != nil {
		return m.encryptWithMetadataFunc(encoded, metadata)
	}
	return crypto.NewEncryptedData([]byte("ciphertext_with_metadata"), metadata), nil
}

func newTestService() (*Service, *mockEncoder, *mockEncryptor) {
	encoder := &mockEncoder{maxSlots: 8192}
	encryptor := &mockEncryptor{}

	keyConfig := crypto.NewKeyConfig("test_hash", []int{1, 2, 4})
	publicKey := crypto.NewKeyMaterial(crypto.KeyTypePublic, []byte("public_key"), keyConfig)
	secretKey := crypto.NewKeyMaterial(crypto.KeyTypeSecret, []byte("secret_key"), keyConfig)
	relinKey := crypto.NewKeyMaterial(crypto.KeyTypeRelinearization, []byte("relin_key"), keyConfig)

	keySet := crypto.NewKeySet("test_hash")
	keySet.SetPublicKey(publicKey)
	keySet.SetSecretKey(secretKey)
	keySet.SetRelinearizationKey(relinKey)

	for _, pos := range []int{1, 2, 4} {
		rotKey := crypto.NewKeyMaterial(crypto.KeyTypeRotation, []byte("rot_key"), keyConfig)
		keySet.AddRotationKey(pos, rotKey)
	}

	params := &crypto.Parameters{}

	service := NewService(encoder, encryptor, keySet, params)
	return service, encoder, encryptor
}

func TestNewService(t *testing.T) {
	service, _, _ := newTestService()

	if service == nil {
		t.Fatal("expected service to be created, got nil")
	}

	if service.encoder == nil {
		t.Error("expected encoder to be set")
	}

	if service.encryptor == nil {
		t.Error("expected encryptor to be set")
	}

	if service.keySet == nil {
		t.Error("expected keySet to be set")
	}

	if service.params == nil {
		t.Error("expected params to be set")
	}
}

func TestService_PrepareRequest(t *testing.T) {
	tests := []struct {
		name          string
		operationType operation.OperationType
		datasets      []*dataset.Dataset
		singleValues  []float64
		parameters    map[string]interface{}
		requestID     string
		wantError     bool
		errorMsg      string
	}{
		{
			name:          "successful request with dataset",
			operationType: operation.OperationVariance,
			datasets: func() []*dataset.Dataset {
				ds, _ := dataset.NewDataset([]float64{1.0, 2.0, 3.0})
				return []*dataset.Dataset{ds}
			}(),
			singleValues: []float64{},
			parameters:   map[string]interface{}{"test_param": "value"},
			requestID:    "test-123",
			wantError:    false,
		},
		{
			name:          "successful request with single values",
			operationType: operation.OperationVariance,
			datasets:      []*dataset.Dataset{},
			singleValues:  []float64{5.0, 10.0},
			parameters:    map[string]interface{}{},
			requestID:     "test-456",
			wantError:     false,
		},
		{
			name:          "successful request with both datasets and single values",
			operationType: operation.OperationVariance,
			datasets: func() []*dataset.Dataset {
				ds, _ := dataset.NewDataset([]float64{1.0, 2.0, 3.0})
				return []*dataset.Dataset{ds}
			}(),
			singleValues: []float64{5.0},
			parameters:   map[string]interface{}{"alpha": 0.05},
			requestID:    "test-789",
			wantError:    false,
		},
		{
			name:          "error encoding dataset",
			operationType: operation.OperationVariance,
			datasets: func() []*dataset.Dataset {
				ds, _ := dataset.NewDataset([]float64{1.0, 2.0, 3.0})
				return []*dataset.Dataset{ds}
			}(),
			singleValues: []float64{},
			parameters:   map[string]interface{}{},
			requestID:    "test-error",
			wantError:    true,
			errorMsg:     "failed to encrypt dataset 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, encoder, _ := newTestService()

			if tt.wantError {
				encoder.encodeFunc = func(values []float64) ([]byte, error) {
					return nil, errors.New("encoding error")
				}
			}

			result, err := service.PrepareRequest(
				tt.operationType,
				tt.datasets,
				tt.singleValues,
				tt.parameters,
				tt.requestID,
			)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected result to be non-nil")
			}

			if result.OperationType() != tt.operationType {
				t.Errorf("expected operation type %s, got %s", tt.operationType, result.OperationType())
			}

			if result.RequestID() != tt.requestID {
				t.Errorf("expected request ID %s, got %s", tt.requestID, result.RequestID())
			}

			expectedInputs := len(tt.datasets) + len(tt.singleValues)
			if len(result.EncryptedInputs()) != expectedInputs {
				t.Errorf("expected %d encrypted inputs, got %d", expectedInputs, len(result.EncryptedInputs()))
			}

			for key, value := range tt.parameters {
				paramValue, exists := result.GetParameter(key)
				if !exists {
					t.Errorf("expected parameter %s to be set", key)
				}
				if paramValue != value {
					t.Errorf("expected parameter %s to be %v, got %v", key, value, paramValue)
				}
			}
		})
	}
}

func TestService_EncryptDataset(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		wantError bool
		errorMsg  string
	}{
		{
			name:      "successful encryption",
			values:    []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			wantError: false,
		},
		{
			name:      "encoding error",
			values:    []float64{1.0, 2.0, 3.0},
			wantError: true,
			errorMsg:  "encoding failed",
		},
		{
			name:      "encryption error",
			values:    []float64{1.0, 2.0, 3.0},
			wantError: true,
			errorMsg:  "encryption failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, encoder, encryptor := newTestService()

			if tt.wantError {
				if contains(tt.errorMsg, "encoding") {
					encoder.encodeFunc = func(values []float64) ([]byte, error) {
						return nil, errors.New("encoding error")
					}
				} else if contains(tt.errorMsg, "encryption") {
					encryptor.encryptFunc = func(encoded []byte) (*crypto.EncryptedData, error) {
						return nil, errors.New("encryption error")
					}
				}
			}

			ds, err := dataset.NewDataset(tt.values)
			if err != nil {
				t.Fatalf("failed to create dataset: %v", err)
			}

			result, err := service.EncryptDataset(ds)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected result to be non-nil")
			}
		})
	}
}

func TestService_EncryptValue(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		wantError bool
		errorMsg  string
	}{
		{
			name:      "successful encryption",
			value:     42.5,
			wantError: false,
		},
		{
			name:      "encoding error",
			value:     10.0,
			wantError: true,
			errorMsg:  "encoding failed",
		},
		{
			name:      "encryption error",
			value:     20.0,
			wantError: true,
			errorMsg:  "encryption failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, encoder, encryptor := newTestService()

			if tt.wantError {
				if contains(tt.errorMsg, "encoding") {
					encoder.encodeSingleFunc = func(value float64, slots int) ([]byte, error) {
						return nil, errors.New("encoding error")
					}
				} else if contains(tt.errorMsg, "encryption") {
					encryptor.encryptFunc = func(encoded []byte) (*crypto.EncryptedData, error) {
						return nil, errors.New("encryption error")
					}
				}
			}

			result, err := service.EncryptValue(tt.value)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("expected result to be non-nil")
			}
		})
	}
}

func TestService_GetPublicKeySet(t *testing.T) {
	service, _, _ := newTestService()

	publicKeySet := service.GetPublicKeySet()

	if publicKeySet == nil {
		t.Fatal("expected public key set to be non-nil")
	}

	if publicKeySet.PublicKey() == nil {
		t.Error("expected public key to be set")
	}

	if publicKeySet.RelinearizationKey() == nil {
		t.Error("expected relinearization key to be set")
	}

	if publicKeySet.HasSecretKey() {
		t.Error("expected secret key to not be included in public key set")
	}

	expectedRotations := []int{1, 2, 4}
	rotations := publicKeySet.RotationPositions()

	if len(rotations) != len(expectedRotations) {
		t.Errorf("expected %d rotation keys, got %d", len(expectedRotations), len(rotations))
	}

	for _, pos := range expectedRotations {
		if _, exists := publicKeySet.RotationKey(pos); !exists {
			t.Errorf("expected rotation key for position %d to exist", pos)
		}
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
