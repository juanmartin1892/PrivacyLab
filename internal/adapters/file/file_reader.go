package file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/juanmartin/privacylab/internal/domain/dataset"
	"github.com/juanmartin/privacylab/internal/domain/errors"
)

type FileReaderAdapter struct{}

func NewFileReaderAdapter() *FileReaderAdapter {
	return &FileReaderAdapter{}
}

func (r *FileReaderAdapter) ReadFromFile(filepath string) (*dataset.Dataset, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filepath, err)
	}
	defer file.Close()

	values, err := parseValues(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filepath, err)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("%w: no valid values found in file %s", errors.ErrEmptyDataset, filepath)
	}

	return dataset.NewDataset(values)
}

func (r *FileReaderAdapter) ReadFromBytes(data []byte) (*dataset.Dataset, error) {
	reader := bytes.NewReader(data)
	values, err := parseValues(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("%w: no valid values found in data", errors.ErrEmptyDataset)
	}

	return dataset.NewDataset(values)
}

func parseValues(reader interface{ Read([]byte) (int, error) }) ([]float64, error) {
	scanner := bufio.NewScanner(reader)
	var values []float64
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value on line %d: %s (%w)", lineNumber, line, err)
		}

		values = append(values, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading data: %w", err)
	}

	return values, nil
}
