package lattigo

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type EncryptorAdapter struct {
	encryptor *rlwe.Encryptor
	params    ckks.Parameters
}

func NewEncryptorAdapter(params *crypto.Parameters, publicKey *crypto.KeyMaterial) (*EncryptorAdapter, error) {
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

	pk, err := deserializePublicKey(publicKey.Data())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize public key: %w", err)
	}

	encryptor := rlwe.NewEncryptor(ckksParams, pk)

	return &EncryptorAdapter{
		encryptor: encryptor,
		params:    ckksParams,
	}, nil
}

func (e *EncryptorAdapter) Encrypt(encoded []byte) (*crypto.EncryptedData, error) {
	plaintext, err := deserializePlaintext(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize plaintext: %w", err)
	}

	ciphertext, err := e.encryptor.EncryptNew(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt plaintext: %w", err)
	}

	ctBytes, err := serializeCiphertext(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize ciphertext: %w", err)
	}

	metadata := crypto.NewEncryptionMetadata(
		len(ctBytes),
		ciphertext.Level(),
		ciphertext.Scale.Float64(),
		e.params.MaxSlots(),
		computeParameterHashFromCKKS(e.params),
	)

	return crypto.NewEncryptedData(ctBytes, metadata), nil
}

func (e *EncryptorAdapter) EncryptWithMetadata(encoded []byte, metadata crypto.EncryptionMetadata) (*crypto.EncryptedData, error) {
	plaintext, err := deserializePlaintext(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize plaintext: %w", err)
	}

	ciphertext, err := e.encryptor.EncryptNew(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt plaintext: %w", err)
	}

	ctBytes, err := serializeCiphertext(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize ciphertext: %w", err)
	}

	return crypto.NewEncryptedData(ctBytes, metadata), nil
}

func deserializePublicKey(data []byte) (*rlwe.PublicKey, error) {
	var pk rlwe.PublicKey
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&pk); err != nil {
		return nil, err
	}
	return &pk, nil
}

func serializeCiphertext(ct *rlwe.Ciphertext) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(ct); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func computeParameterHashFromCKKS(params ckks.Parameters) string {
	// TODO: Use a proper hash function if needed
	return fmt.Sprintf("CKKS_n%d_q%v_p%v_s%d",
		params.LogN(),
		params.LogQ(),
		params.LogP(),
		params.LogDefaultScale(),
	)
}
