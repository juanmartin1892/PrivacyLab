package ports

// Encoder defines the contract for encoding plaintext values into representations
// suitable for homomorphic encryption.
type Encoder interface {
	// Encode encodes a slice of values into an encoded plaintext.
	// Values are packed into SIMD slots for efficient batch processing.
	Encode(values []float64) ([]byte, error)

	// EncodeSingle encodes a single value, replicating it across all slots.
	// Used for operations where the same value needs to be applied to all slots.
	EncodeSingle(value float64, slots int) ([]byte, error)

	// Decode decodes an encoded plaintext back to a slice of values.
	Decode(encoded []byte, slots int) ([]float64, error)

	// MaxSlots returns the maximum number of SIMD slots available.
	MaxSlots() int
}
