# File Adapters

This directory contains adapter implementations for reading datasets from external sources.

## Overview

The file adapter implements the `ports.DataReader` interface, providing functionality to load numerical datasets from text files and byte arrays.

## FileReaderAdapter

Implements `ports.DataReader` for reading datasets from files.

**Features:**
- Read datasets from text files (one value per line)
- Read datasets from byte arrays
- Automatic parsing and validation of numeric values
- Support for comments (lines starting with `#`)
- Skip empty lines

**File Format:**

The expected file format is simple text with one numeric value per line:

```
# Population data
10.5
20.3
15.7
8.9
12.4
```

**Usage:**

```go
import (
    "github.com/juanmartin/privacylab/internal/adapters/file"
    "github.com/juanmartin/privacylab/internal/domain/dataset"
)

// Create reader
reader := file.NewFileReaderAdapter()

// Read from file
dataset, err := reader.ReadFromFile("population.txt")
if err != nil {
    log.Fatalf("Failed to read dataset: %v", err)
}

fmt.Printf("Loaded %d values\n", dataset.Size())
values := dataset.Values()

// Read from bytes
data := []byte("1.5\n2.7\n3.9\n4.2\n")
dataset2, err := reader.ReadFromBytes(data)
if err != nil {
    log.Fatalf("Failed to parse data: %v", err)
}
```

## Error Handling

The adapter returns descriptive errors for common issues:

- **File not found:** `failed to open file <path>: no such file or directory`
- **Empty file:** `ErrEmptyDataset: no valid values found in file`
- **Invalid value:** `invalid value on line <N>: <value>`
- **Invalid number:** `ErrInvalidDataValue: invalid value at index <N>` (NaN or Inf)

## Example Files

See `examples/point-to-point/population.txt` for a sample dataset file.

## Testing

Example test:

```go
func TestFileReaderAdapter_ReadFromFile(t *testing.T) {
    reader := file.NewFileReaderAdapter()
    
    dataset, err := reader.ReadFromFile("testdata/values.txt")
    
    assert.NoError(t, err)
    assert.NotNil(t, dataset)
    assert.Greater(t, dataset.Size(), 0)
}

func TestFileReaderAdapter_ReadFromBytes(t *testing.T) {
    reader := file.NewFileReaderAdapter()
    data := []byte("1.0\n2.0\n3.0\n")
    
    dataset, err := reader.ReadFromBytes(data)
    
    assert.NoError(t, err)
    assert.Equal(t, 3, dataset.Size())
}
```

## Supported Formats

Currently supports:
- Plain text files with one float64 value per line
- Comments starting with `#`
- Empty lines (ignored)

Future enhancements could include:
- CSV format support
- JSON format support
- Binary format support
- Streaming for large files
