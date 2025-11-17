package lattigo

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/juanmartin/privacylab/internal/domain/computation"
	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/operation"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

type EvaluatorAdapter struct {
	evaluator *ckks.Evaluator
	params    ckks.Parameters
	workerID  string
}

func NewEvaluatorAdapter(
	params *crypto.Parameters,
	keySet *crypto.KeySet,
	workerID string,
) (*EvaluatorAdapter, error) {
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

	rlk, err := deserializeRelinearizationKey(keySet.RelinearizationKey().Data())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize relinearization key: %w", err)
	}

	rotationPositions := keySet.RotationPositions()
	galoisKeys := make([]*rlwe.GaloisKey, 0, len(rotationPositions))
	for _, pos := range rotationPositions {
		rotKey, exists := keySet.RotationKey(pos)
		if !exists {
			continue
		}
		gk, err := deserializeGaloisKey(rotKey.Data())
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize rotation key for position %d: %w", pos, err)
		}
		galoisKeys = append(galoisKeys, gk)
	}

	evk := rlwe.NewMemEvaluationKeySet(rlk, galoisKeys...)
	evaluator := ckks.NewEvaluator(ckksParams, evk)

	return &EvaluatorAdapter{
		evaluator: evaluator,
		params:    ckksParams,
		workerID:  workerID,
	}, nil
}

func (e *EvaluatorAdapter) EvaluateOperation(
	op operation.StatisticalOperation,
	request *computation.Request,
) (*computation.Result, error) {
	startTime := time.Now()

	var result *crypto.EncryptedData
	var multiplicationsUsed, rotationsUsed int
	var err error

	switch op.Type() {
	case operation.VarianceType:
		result, multiplicationsUsed, rotationsUsed, err = e.evaluateVariance(request)
	default:
		return nil, fmt.Errorf("unsupported operation type: %s", op.Type())
	}

	if err != nil {
		return nil, err
	}

	computationTime := time.Since(startTime)

	metadata := computation.NewResultMetadata(
		computationTime,
		multiplicationsUsed,
		rotationsUsed,
		result.Level(),
		result.Level(),
		result.Scale(),
		e.workerID,
	)

	return computation.NewResult(result, op.Type(), request.RequestID(), metadata), nil
}

func (e *EvaluatorAdapter) evaluateVariance(request *computation.Request) (*crypto.EncryptedData, int, int, error) {
	inputs := request.EncryptedInputs()
	if len(inputs) < 3 {
		return nil, 0, 0, fmt.Errorf("variance requires 3 inputs: population, reference value, and n")
	}

	encryptedPopulation := inputs[0]
	encryptedX := inputs[1]

	dataSizeParam, exists := request.GetParameter("data_size")
	if !exists {
		return nil, 0, 0, fmt.Errorf("data_size parameter is required for variance computation")
	}
	dataSize, ok := dataSizeParam.(int)
	if !ok {
		return nil, 0, 0, fmt.Errorf("data_size parameter must be an integer")
	}

	ctPopulation, err := deserializeCiphertext(encryptedPopulation.Ciphertext())
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to deserialize population ciphertext: %w", err)
	}

	ctX, err := deserializeCiphertext(encryptedX.Ciphertext())
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to deserialize reference value ciphertext: %w", err)
	}

	diff, err := e.evaluator.SubNew(ctPopulation, ctX)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to compute difference: %w", err)
	}

	squaredDiff, err := e.evaluator.MulRelinNew(diff, diff)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to compute squared difference: %w", err)
	}

	if err := e.evaluator.Rescale(squaredDiff, squaredDiff); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to rescale squared difference: %w", err)
	}

	sumResult, rotations := e.sumSlotsInternal(squaredDiff, dataSize)

	invN := 1.0 / float64(dataSize)
	variance, err := e.evaluator.MulNew(sumResult, invN)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to compute final variance: %w", err)
	}

	if err := e.evaluator.Rescale(variance, variance); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to rescale final result: %w", err)
	}

	varianceBytes, err := serializeCiphertext(variance)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to serialize result: %w", err)
	}

	metadata := crypto.NewEncryptionMetadata(
		len(varianceBytes),
		variance.Level(),
		variance.Scale.Float64(),
		e.params.MaxSlots(),
		computeParameterHashFromCKKS(e.params),
	)

	return crypto.NewEncryptedData(varianceBytes, metadata), 2, rotations, nil
}

func (e *EvaluatorAdapter) Add(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	ctA, err := deserializeCiphertext(a.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize first ciphertext: %w", err)
	}

	ctB, err := deserializeCiphertext(b.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize second ciphertext: %w", err)
	}

	result, err := e.evaluator.AddNew(ctA, ctB)
	if err != nil {
		return nil, fmt.Errorf("failed to perform addition: %w", err)
	}

	return e.ciphertextToEncryptedData(result)
}

func (e *EvaluatorAdapter) Sub(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	ctA, err := deserializeCiphertext(a.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize first ciphertext: %w", err)
	}

	ctB, err := deserializeCiphertext(b.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize second ciphertext: %w", err)
	}

	result, err := e.evaluator.SubNew(ctA, ctB)
	if err != nil {
		return nil, fmt.Errorf("failed to perform subtraction: %w", err)
	}

	return e.ciphertextToEncryptedData(result)
}

func (e *EvaluatorAdapter) Mul(a, b *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	ctA, err := deserializeCiphertext(a.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize first ciphertext: %w", err)
	}

	ctB, err := deserializeCiphertext(b.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize second ciphertext: %w", err)
	}

	result, err := e.evaluator.MulRelinNew(ctA, ctB)
	if err != nil {
		return nil, fmt.Errorf("failed to perform multiplication: %w", err)
	}

	return e.ciphertextToEncryptedData(result)
}

func (e *EvaluatorAdapter) MulScalar(ct *crypto.EncryptedData, scalar float64) (*crypto.EncryptedData, error) {
	ciphertext, err := deserializeCiphertext(ct.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ciphertext: %w", err)
	}

	result, err := e.evaluator.MulNew(ciphertext, scalar)
	if err != nil {
		return nil, fmt.Errorf("failed to perform scalar multiplication: %w", err)
	}

	return e.ciphertextToEncryptedData(result)
}

func (e *EvaluatorAdapter) Rotate(ct *crypto.EncryptedData, positions int) (*crypto.EncryptedData, error) {
	ciphertext, err := deserializeCiphertext(ct.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ciphertext: %w", err)
	}

	result, err := e.evaluator.RotateNew(ciphertext, positions)
	if err != nil {
		return nil, fmt.Errorf("failed to perform rotation: %w", err)
	}

	return e.ciphertextToEncryptedData(result)
}

func (e *EvaluatorAdapter) Rescale(ct *crypto.EncryptedData) (*crypto.EncryptedData, error) {
	ciphertext, err := deserializeCiphertext(ct.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ciphertext: %w", err)
	}

	result := ciphertext.CopyNew()
	if err := e.evaluator.Rescale(result, result); err != nil {
		return nil, fmt.Errorf("failed to perform rescaling: %w", err)
	}

	return e.ciphertextToEncryptedData(result)
}

func (e *EvaluatorAdapter) SumSlots(ct *crypto.EncryptedData, n int) (*crypto.EncryptedData, error) {
	ciphertext, err := deserializeCiphertext(ct.Ciphertext())
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize ciphertext: %w", err)
	}

	result, _ := e.sumSlotsInternal(ciphertext, n)

	return e.ciphertextToEncryptedData(result)
}

func (e *EvaluatorAdapter) sumSlotsInternal(ct *rlwe.Ciphertext, n int) (*rlwe.Ciphertext, int) {
	result := ct.CopyNew()
	rotationsUsed := 0

	for i := 1; i < n; i++ {
		rotated, err := e.evaluator.RotateNew(ct, i)
		if err != nil {
			log.Printf("Error rotating ciphertext at position %d: %v", i, err)
			continue
		}

		if err := e.evaluator.Add(result, rotated, result); err != nil {
			log.Printf("Error adding rotated ciphertext: %v", err)
			continue
		}

		rotationsUsed++
	}

	return result, rotationsUsed
}

func (e *EvaluatorAdapter) ciphertextToEncryptedData(ct *rlwe.Ciphertext) (*crypto.EncryptedData, error) {
	ctBytes, err := serializeCiphertext(ct)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize ciphertext: %w", err)
	}

	metadata := crypto.NewEncryptionMetadata(
		len(ctBytes),
		ct.Level(),
		ct.Scale.Float64(),
		e.params.MaxSlots(),
		computeParameterHashFromCKKS(e.params),
	)

	return crypto.NewEncryptedData(ctBytes, metadata), nil
}

func deserializeRelinearizationKey(data []byte) (*rlwe.RelinearizationKey, error) {
	var rlk rlwe.RelinearizationKey
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&rlk); err != nil {
		return nil, err
	}
	return &rlk, nil
}

func deserializeGaloisKey(data []byte) (*rlwe.GaloisKey, error) {
	var gk rlwe.GaloisKey
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&gk); err != nil {
		return nil, err
	}
	return &gk, nil
}
