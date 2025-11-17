package lattigo

import (
	"testing"

	"github.com/juanmartin/privacylab/internal/domain/operation"
)

func TestNewRotationMapper(t *testing.T) {
	mapper := NewRotationMapper()

	if mapper == nil {
		t.Fatal("expected mapper to be created")
	}

	if mapper.mappings == nil {
		t.Fatal("expected mappings to be initialized")
	}
}

func TestRotationMapper_GetRequiredRotations_Variance(t *testing.T) {
	mapper := NewRotationMapper()

	tests := []struct {
		name     string
		dataSize int
		want     []int
	}{
		{
			name:     "empty dataset",
			dataSize: 0,
			want:     []int{},
		},
		{
			name:     "single element",
			dataSize: 1,
			want:     []int{},
		},
		{
			name:     "small dataset",
			dataSize: 5,
			want:     []int{1, 2, 3, 4},
		},
		{
			name:     "larger dataset",
			dataSize: 10,
			want:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapper.GetRequiredRotations(operation.VarianceType, tt.dataSize)

			if len(got) != len(tt.want) {
				t.Errorf("got %d rotations, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("rotation[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestRotationMapper_GetRequiredRotations_UnknownOperation(t *testing.T) {
	mapper := NewRotationMapper()

	unknownOp := operation.OperationType("unknown_operation")
	rotations := mapper.GetRequiredRotations(unknownOp, 10)

	if rotations != nil {
		t.Errorf("expected nil for unknown operation, got %v", rotations)
	}
}

func TestRotationMapper_Register(t *testing.T) {
	mapper := NewRotationMapper()

	customOp := operation.OperationType("custom_mean")
	customStrategy := func(dataSize int) []int {
		return []int{1, 2}
	}

	mapper.Register(customOp, customStrategy)

	rotations := mapper.GetRequiredRotations(customOp, 100)
	expected := []int{1, 2}

	if len(rotations) != len(expected) {
		t.Errorf("got %d rotations, want %d", len(rotations), len(expected))
		return
	}

	for i := range rotations {
		if rotations[i] != expected[i] {
			t.Errorf("rotation[%d] = %d, want %d", i, rotations[i], expected[i])
		}
	}
}

func TestVarianceRotationStrategy(t *testing.T) {
	tests := []struct {
		name     string
		dataSize int
		want     []int
	}{
		{
			name:     "zero size",
			dataSize: 0,
			want:     []int{},
		},
		{
			name:     "size one",
			dataSize: 1,
			want:     []int{},
		},
		{
			name:     "size five",
			dataSize: 5,
			want:     []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VarianceRotationStrategy(tt.dataSize)

			if len(got) != len(tt.want) {
				t.Errorf("got %d rotations, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("rotation[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}
