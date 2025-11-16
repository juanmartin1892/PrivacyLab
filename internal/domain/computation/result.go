package computation

import (
	"time"

	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/operation"
)

// Result represents the result of a homomorphic computation performed by the worker.
type Result struct {
	encryptedResult *crypto.EncryptedData
	operationType   operation.OperationType
	requestID       string
	metadata        ResultMetadata
}

// ResultMetadata contains information about the computation performed.
type ResultMetadata struct {
	computationTime     time.Duration
	multiplicationsUsed int
	rotationsUsed       int
	levelsConsumed      int
	finalLevel          int
	finalScale          float64
	workerID            string
	timestamp           time.Time
}

// NewResult creates a new computation result.
func NewResult(
	encryptedResult *crypto.EncryptedData,
	operationType operation.OperationType,
	requestID string,
	metadata ResultMetadata,
) *Result {
	return &Result{
		encryptedResult: encryptedResult,
		operationType:   operationType,
		requestID:       requestID,
		metadata:        metadata,
	}
}

// EncryptedResult returns the encrypted computation result.
func (r *Result) EncryptedResult() *crypto.EncryptedData {
	return r.encryptedResult
}

// OperationType returns the type of operation that was performed.
func (r *Result) OperationType() operation.OperationType {
	return r.operationType
}

// RequestID returns the request identifier this result corresponds to.
func (r *Result) RequestID() string {
	return r.requestID
}

// Metadata returns the computation metadata.
func (r *Result) Metadata() ResultMetadata {
	return r.metadata
}

// ComputationTime returns the time taken to perform the computation.
func (r *Result) ComputationTime() time.Duration {
	return r.metadata.computationTime
}

// NewResultMetadata creates a new result metadata object.
func NewResultMetadata(
	computationTime time.Duration,
	multiplicationsUsed int,
	rotationsUsed int,
	levelsConsumed int,
	finalLevel int,
	finalScale float64,
	workerID string,
) ResultMetadata {
	return ResultMetadata{
		computationTime:     computationTime,
		multiplicationsUsed: multiplicationsUsed,
		rotationsUsed:       rotationsUsed,
		levelsConsumed:      levelsConsumed,
		finalLevel:          finalLevel,
		finalScale:          finalScale,
		workerID:            workerID,
		timestamp:           time.Now(),
	}
}

// MultiplicationsUsed returns the number of homomorphic multiplications performed.
func (m ResultMetadata) MultiplicationsUsed() int {
	return m.multiplicationsUsed
}

// RotationsUsed returns the number of slot rotations performed.
func (m ResultMetadata) RotationsUsed() int {
	return m.rotationsUsed
}

// LevelsConsumed returns the number of multiplication levels consumed.
func (m ResultMetadata) LevelsConsumed() int {
	return m.levelsConsumed
}

// FinalLevel returns the final multiplication level of the result.
func (m ResultMetadata) FinalLevel() int {
	return m.finalLevel
}

// FinalScale returns the final scaling factor of the result.
func (m ResultMetadata) FinalScale() float64 {
	return m.finalScale
}

// WorkerID returns the identifier of the worker that performed the computation.
func (m ResultMetadata) WorkerID() string {
	return m.workerID
}

// Timestamp returns when the computation was completed.
func (m ResultMetadata) Timestamp() time.Time {
	return m.timestamp
}
