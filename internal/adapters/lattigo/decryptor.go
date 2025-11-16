package lattigo

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type DecryptorAdapter struct {
	decryptor *rlwe.Decryptor
	encoder   *ckks.Encoder
	params    ckks.Parameters
}

func NewDecryptorAdapter(params *crypto.Parameters, secretKey *crypto.KeyMaterial) (*DecryptorAdapter, error) {
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

	sk, err := deserializeSecretKey(secretKey.Data())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize secret key: %w", err)
	}

	decryptor := rlwe.NewDecryptor(ckksParams, sk)
	encoder := ckks.NewEncoder(ckksParams)

	return &DecryptorAdapter{
		decryptor: decryptor,
		encoder:   encoder,
		params:    ckksParams,
	}, nil
}

func (d *DecryptorAdapter) Decrypt(encrypted *crypto.EncryptedData) ([]byte, error) {
	ciphertext, err := deserializeCiphertext(encrypted.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ciphertext: %w", err)
	}

	plaintext := d.decryptor.DecryptNew(ciphertext)

	return serializePlaintext(plaintext)
}

func (d *DecryptorAdapter) DecryptValue(encrypted *crypto.EncryptedData) (float64, error) {
	ciphertext, err := deserializeCiphertext(encrypted.Ciphertext())
	if err != nil {
		return 0, fmt.Errorf("failed to deserialize ciphertext: %w", err)
	}

	plaintext := d.decryptor.DecryptNew(ciphertext)

	result := make([]complex128, d.params.MaxSlots())
	if err := d.encoder.Decode(plaintext, result); err != nil {
		return 0, fmt.Errorf("failed to decode plaintext: %w", err)
	}

	return real(result[0]), nil
}

func deserializeCiphertext(data []byte) (*rlwe.Ciphertext, error) {
	var ct rlwe.Ciphertext
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&ct); err != nil {
		return nil, err
	}
	return &ct, nil
}
