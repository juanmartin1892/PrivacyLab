package lattigo

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type EncoderAdapter struct {
	encoder *ckks.Encoder
	params  ckks.Parameters
}

func NewEncoderAdapter(params *crypto.Parameters) (*EncoderAdapter, error) {
	if params.Scheme() != crypto.SchemeCKKS {
		return nil, fmt.Errorf("only CKKS scheme is supported, got %s", params.Scheme())
	}

	ckksParams, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            params.LogN(),
		LogQ:            params.LogQ(),
		LogP:            params.LogP(),
		LogDefaultScale: params.LogDefaultScale(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create CKKS parameters: %w", err)
	}

	encoder := ckks.NewEncoder(ckksParams)

	return &EncoderAdapter{
		encoder: encoder,
		params:  ckksParams,
	}, nil
}

func (e *EncoderAdapter) Encode(values []float64) ([]byte, error) {
	if len(values) > e.params.MaxSlots() {
		return nil, fmt.Errorf("values exceed maximum slots: got %d, max %d", len(values), e.params.MaxSlots())
	}

	plaintext := ckks.NewPlaintext(e.params, e.params.MaxLevel())
	if err := e.encoder.Encode(values, plaintext); err != nil {
		return nil, fmt.Errorf("failed to encode values: %w", err)
	}

	return serializePlaintext(plaintext)
}

func (e *EncoderAdapter) EncodeSingle(value float64, slots int) ([]byte, error) {
	if slots > e.params.MaxSlots() {
		slots = e.params.MaxSlots()
	}

	values := make([]float64, slots)
	for i := range values {
		values[i] = value
	}

	plaintext := ckks.NewPlaintext(e.params, e.params.MaxLevel())
	if err := e.encoder.Encode(values, plaintext); err != nil {
		return nil, fmt.Errorf("failed to encode single value: %w", err)
	}

	return serializePlaintext(plaintext)
}

func (e *EncoderAdapter) Decode(encoded []byte, slots int) ([]float64, error) {
	plaintext, err := deserializePlaintext(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize plaintext: %w", err)
	}

	if slots > e.params.MaxSlots() {
		slots = e.params.MaxSlots()
	}

	result := make([]complex128, slots)
	if err := e.encoder.Decode(plaintext, result); err != nil {
		return nil, fmt.Errorf("failed to decode plaintext: %w", err)
	}

	values := make([]float64, slots)
	for i := range result {
		values[i] = real(result[i])
	}

	return values, nil
}

func (e *EncoderAdapter) MaxSlots() int {
	return e.params.MaxSlots()
}

func serializePlaintext(pt *rlwe.Plaintext) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(pt); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserializePlaintext(data []byte) (*rlwe.Plaintext, error) {
	var pt rlwe.Plaintext
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&pt); err != nil {
		return nil, err
	}
	return &pt, nil
}
