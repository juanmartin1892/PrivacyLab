# Dataset Domain

This package defines the Dataset entity, which represents a collection of numerical values to be processed using homomorphic encryption.

## Purpose

The Dataset entity encapsulates raw numerical data and provides:
- **Immutability**: Values are copied defensively to prevent external modifications
- **Validation**: All values are checked for validity (no NaN or Inf values)
- **Type safety**: Provides a strongly-typed interface for numerical collections

## Key Concepts

A Dataset is a value object in the domain that represents the input data for statistical operations. It enforces invariants at construction time, ensuring that:
- The dataset is never empty
- All values are valid finite floating-point numbers
- External code cannot modify the internal state

## Usage

```go
import "github.com/juanmartin/privacyLab/internal/domain/dataset"

// Create a new dataset
values := []float64{1.5, 2.3, 4.7, 3.1}
ds, err := dataset.NewDataset(values)
if err != nil {
    // Handle validation error
}

// Access values safely
size := ds.Size()
allValues := ds.Values() // Returns a copy
```

## Design Decisions

- Uses defensive copying to maintain immutability
- Validates data at construction time (fail-fast principle)
- Does not depend on any infrastructure libraries
