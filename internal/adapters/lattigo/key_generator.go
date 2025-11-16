package lattigo

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type KeyGeneratorAdapter struct {
	params ckks.Parameters
	keygen *rlwe.KeyGenerator
}

func NewKeyGeneratorAdapter(params *crypto.Parameters) (*KeyGeneratorAdapter, error) {
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

	keygen := rlwe.NewKeyGenerator(ckksParams)

	return &KeyGeneratorAdapter{
		params: ckksParams,
		keygen: keygen,
	}, nil
}

func (kg *KeyGeneratorAdapter) GenerateKeys(params *crypto.Parameters, rotations []int) (*crypto.KeySet, error) {
	paramHash := computeParameterHash(params)
	keySet := crypto.NewKeySet(paramHash)

	secretKey, err := kg.GenerateSecretKey(params)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %w", err)
	}
	keySet.SetSecretKey(secretKey)

	publicKey, err := kg.GeneratePublicKey(params, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public key: %w", err)
	}
	keySet.SetPublicKey(publicKey)

	relinKey, err := kg.GenerateRelinearizationKey(params, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate relinearization key: %w", err)
	}
	keySet.SetRelinearizationKey(relinKey)

	if len(rotations) > 0 {
		rotationKeys, err := kg.GenerateRotationKeys(params, secretKey, rotations)
		if err != nil {
			return nil, fmt.Errorf("failed to generate rotation keys: %w", err)
		}
		for pos, key := range rotationKeys {
			keySet.AddRotationKey(pos, key)
		}
	}

	return keySet, nil
}

func (kg *KeyGeneratorAdapter) GenerateSecretKey(params *crypto.Parameters) (*crypto.KeyMaterial, error) {
	sk := kg.keygen.GenSecretKeyNew()

	skBytes, err := serializeSecretKey(sk)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize secret key: %w", err)
	}

	paramHash := computeParameterHash(params)
	config := crypto.NewKeyConfig(paramHash, nil)

	return crypto.NewKeyMaterial(crypto.KeyTypeSecret, skBytes, config), nil
}

func (kg *KeyGeneratorAdapter) GeneratePublicKey(params *crypto.Parameters, secretKey *crypto.KeyMaterial) (*crypto.KeyMaterial, error) {
	sk, err := deserializeSecretKey(secretKey.Data())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize secret key: %w", err)
	}

	pk := kg.keygen.GenPublicKeyNew(sk)

	pkBytes, err := serializePublicKey(pk)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize public key: %w", err)
	}

	paramHash := computeParameterHash(params)
	config := crypto.NewKeyConfig(paramHash, nil)

	return crypto.NewKeyMaterial(crypto.KeyTypePublic, pkBytes, config), nil
}

func (kg *KeyGeneratorAdapter) GenerateRelinearizationKey(params *crypto.Parameters, secretKey *crypto.KeyMaterial) (*crypto.KeyMaterial, error) {
	sk, err := deserializeSecretKey(secretKey.Data())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize secret key: %w", err)
	}

	rlk := kg.keygen.GenRelinearizationKeyNew(sk)

	rlkBytes, err := serializeRelinearizationKey(rlk)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize relinearization key: %w", err)
	}

	paramHash := computeParameterHash(params)
	config := crypto.NewKeyConfig(paramHash, nil)

	return crypto.NewKeyMaterial(crypto.KeyTypeRelinearization, rlkBytes, config), nil
}

func (kg *KeyGeneratorAdapter) GenerateRotationKeys(params *crypto.Parameters, secretKey *crypto.KeyMaterial, positions []int) (map[int]*crypto.KeyMaterial, error) {
	sk, err := deserializeSecretKey(secretKey.Data())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize secret key: %w", err)
	}

	galoisElements := make([]uint64, len(positions))
	for i, pos := range positions {
		galoisElements[i] = kg.params.GaloisElement(pos)
	}

	rotKeys := kg.keygen.GenGaloisKeysNew(galoisElements, sk)

	result := make(map[int]*crypto.KeyMaterial)
	paramHash := computeParameterHash(params)

	for i, pos := range positions {
		rotKey := rotKeys[i]

		rotKeyBytes, err := serializeGaloisKey(rotKey)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize rotation key for position %d: %w", pos, err)
		}

		config := crypto.NewKeyConfig(paramHash, []int{pos})
		result[pos] = crypto.NewKeyMaterial(crypto.KeyTypeRotation, rotKeyBytes, config)
	}

	return result, nil
}

func serializeSecretKey(sk *rlwe.SecretKey) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(sk); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserializeSecretKey(data []byte) (*rlwe.SecretKey, error) {
	var sk rlwe.SecretKey
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&sk); err != nil {
		return nil, err
	}
	return &sk, nil
}

func serializePublicKey(pk *rlwe.PublicKey) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(pk); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func serializeRelinearizationKey(rlk *rlwe.RelinearizationKey) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(rlk); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func serializeGaloisKey(gk *rlwe.GaloisKey) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(gk); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func computeParameterHash(params *crypto.Parameters) string {
	return fmt.Sprintf("%s_n%d_q%v_p%v_s%d",
		params.Scheme(),
		params.LogN(),
		params.LogQ(),
		params.LogP(),
		params.LogDefaultScale(),
	)
}
