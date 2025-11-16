package crypto

// KeyType represents the type of cryptographic key.
type KeyType string

const (
	KeyTypeSecret          KeyType = "secret"
	KeyTypePublic          KeyType = "public"
	KeyTypeRelinearization KeyType = "relinearization"
	KeyTypeRotation        KeyType = "rotation"
	KeyTypeEvaluation      KeyType = "evaluation"
)

// KeyMaterial represents cryptographic keys used in homomorphic encryption.
// Keys are stored as raw bytes to maintain infrastructure independence.
type KeyMaterial struct {
	keyType KeyType
	data    []byte
	config  KeyConfig
}

// KeyConfig contains configuration information for a key.
type KeyConfig struct {
	parameterHash string
	rotations     []int
}

// NewKeyMaterial creates a new key material object.
func NewKeyMaterial(keyType KeyType, data []byte, config KeyConfig) *KeyMaterial {
	return &KeyMaterial{
		keyType: keyType,
		data:    copyBytes(data),
		config:  config,
	}
}

// Type returns the key type.
func (k *KeyMaterial) Type() KeyType {
	return k.keyType
}

// Data returns a copy of the key bytes.
func (k *KeyMaterial) Data() []byte {
	return copyBytes(k.data)
}

// Config returns the key configuration.
func (k *KeyMaterial) Config() KeyConfig {
	return k.config
}

// IsCompatibleWith checks if this key is compatible with the given parameters.
func (k *KeyMaterial) IsCompatibleWith(parameterHash string) bool {
	return k.config.parameterHash == parameterHash
}

// NewKeyConfig creates a new key configuration.
func NewKeyConfig(parameterHash string, rotations []int) KeyConfig {
	rotationsCopy := make([]int, len(rotations))
	copy(rotationsCopy, rotations)
	return KeyConfig{
		parameterHash: parameterHash,
		rotations:     rotationsCopy,
	}
}

// KeySet represents a collection of keys needed for encryption operations.
type KeySet struct {
	public          *KeyMaterial
	secret          *KeyMaterial
	relinearization *KeyMaterial
	rotation        map[int]*KeyMaterial
	parameterHash   string
}

// NewKeySet creates a new key set.
func NewKeySet(parameterHash string) *KeySet {
	return &KeySet{
		rotation:      make(map[int]*KeyMaterial),
		parameterHash: parameterHash,
	}
}

// SetPublicKey sets the public key.
func (ks *KeySet) SetPublicKey(key *KeyMaterial) {
	ks.public = key
}

// SetSecretKey sets the secret key.
func (ks *KeySet) SetSecretKey(key *KeyMaterial) {
	ks.secret = key
}

// SetRelinearizationKey sets the relinearization key.
func (ks *KeySet) SetRelinearizationKey(key *KeyMaterial) {
	ks.relinearization = key
}

// AddRotationKey adds a rotation key for a specific position.
func (ks *KeySet) AddRotationKey(position int, key *KeyMaterial) {
	ks.rotation[position] = key
}

// PublicKey returns the public key.
func (ks *KeySet) PublicKey() *KeyMaterial {
	return ks.public
}

// SecretKey returns the secret key.
func (ks *KeySet) SecretKey() *KeyMaterial {
	return ks.secret
}

// RelinearizationKey returns the relinearization key.
func (ks *KeySet) RelinearizationKey() *KeyMaterial {
	return ks.relinearization
}

// RotationKey returns the rotation key for a specific position.
func (ks *KeySet) RotationKey(position int) (*KeyMaterial, bool) {
	key, exists := ks.rotation[position]
	return key, exists
}

// HasSecretKey checks if the key set contains a secret key.
func (ks *KeySet) HasSecretKey() bool {
	return ks.secret != nil
}

// RotationPositions returns all available rotation positions.
func (ks *KeySet) RotationPositions() []int {
	positions := make([]int, 0, len(ks.rotation))
	for pos := range ks.rotation {
		positions = append(positions, pos)
	}
	return positions
}
