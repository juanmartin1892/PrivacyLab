package ports

import (
	"github.com/juanmartin/privacylab/internal/domain/dataset"
)

// DataReader defines the contract for reading datasets from external sources.
type DataReader interface {
	// ReadFromFile reads a dataset from a file.
	// The file format should contain one numeric value per line.
	ReadFromFile(filepath string) (*dataset.Dataset, error)

	// ReadFromBytes reads a dataset from raw bytes.
	ReadFromBytes(data []byte) (*dataset.Dataset, error)
}
