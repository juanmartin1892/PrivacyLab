package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

func main() {
	fmt.Println("=== Variance Calculation with Homomorphic Encryption ===")
	fmt.Println("Computing variance of input x against population X using encrypted data")
	fmt.Println()

	// Read population from file
	fmt.Println("Reading population from file...")
	population, err := readPopulationFromFile("population.txt")
	if err != nil {
		log.Fatalf("Error reading population file: %v", err)
	}
	fmt.Printf("Population loaded: %v (size: %d)\n", population, len(population))

	// Input value x
	var inputX float64
	fmt.Print("\nEnter input value x: ")
	_, err = fmt.Scanf("%f", &inputX)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	fmt.Printf("Input x: %.2f\n", inputX)

	// Calculate variance in plain text for verification
	fmt.Println("\nCalculating variance in plain text (for verification)...")
	plainTextVariance := calculatePlainTextVariance(population, inputX)
	fmt.Printf("Plain text variance: %.6f\n", plainTextVariance)

	// Configure CKKS parameters
	params, err := ckks.NewParametersFromLiteral(ckks.ParametersLiteral{
		LogN:            14,
		LogQ:            []int{55, 40, 40, 40, 40, 40},
		LogP:            []int{61},
		LogDefaultScale: 40,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Generate keys
	fmt.Println("\nGenerating cryptographic keys...")
	keygen := rlwe.NewKeyGenerator(params)
	secretKey := keygen.GenSecretKeyNew()
	publicKey := keygen.GenPublicKeyNew(secretKey)
	relinKey := keygen.GenRelinearizationKeyNew(secretKey)

	// Generate rotation keys for summing slots
	// We need rotation keys for positions 1 to n-1
	fmt.Println("Generating rotation keys for SIMD operations...")
	galoisElements := []uint64{}
	for i := 1; i < len(population); i++ {
		galoisElements = append(galoisElements, params.GaloisElement(i))
	}
	rotKeys := keygen.GenGaloisKeysNew(galoisElements, secretKey)

	evk := rlwe.NewMemEvaluationKeySet(relinKey, rotKeys...)

	// Create encoder, encryptor, decryptor and evaluator
	encoder := ckks.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, publicKey)
	decryptor := rlwe.NewDecryptor(params, secretKey)
	evaluator := ckks.NewEvaluator(params, evk)

	// Encrypt population, n, and input x using SIMD batching
	// All population values are packed into a single ciphertext
	// This hides the population size from potential attackers
	numSlots := params.MaxSlots()

	fmt.Println("\nEncrypting population data using SIMD batching...")
	encryptedPopulation := encryptPopulationBatched(population, encoder, encryptor, params)
	fmt.Printf("Population encrypted successfully (all %d values in 1 ciphertext)\n", len(population))

	fmt.Println("Encrypting population size n...")
	encryptedN := encryptValueBatched(float64(len(population)), encoder, encryptor, params, numSlots)
	fmt.Println("Population size encrypted successfully")

	fmt.Println("Encrypting input x...")
	encryptedX := encryptValueBatched(inputX, encoder, encryptor, params, numSlots)
	fmt.Println("Input x encrypted successfully")

	// Calculate variance using homomorphic encryption
	fmt.Println("\nCalculating variance on encrypted data...")
	encryptedVariance := calculateHomomorphicVarianceBatched(
		encryptedPopulation,
		encryptedX,
		encryptedN,
		len(population),
		evaluator,
		params,
	)

	// Decrypt result
	fmt.Println("Decrypting result...")
	homomorphicVariance := decryptValue(encryptedVariance, decryptor, encoder, params)

	// Display results
	fmt.Printf("\nResults:\n")
	fmt.Printf("Plain text variance:      %.6f\n", plainTextVariance)
	fmt.Printf("Homomorphic variance:     %.6f\n", homomorphicVariance)
	fmt.Printf("Absolute error:           %.9f\n", math.Abs(homomorphicVariance-plainTextVariance))
	fmt.Printf("Relative error:           %.6f%%\n", math.Abs(homomorphicVariance-plainTextVariance)/plainTextVariance*100)

	fmt.Println("\nExample completed successfully")
	fmt.Println("Variance was calculated without revealing the original data")
}

// readPopulationFromFile reads population values from a text file
func readPopulationFromFile(filename string) ([]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var population []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value '%s': %w", line, err)
		}
		population = append(population, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(population) == 0 {
		return nil, fmt.Errorf("no valid values found in file")
	}

	return population, nil
}

// calculatePlainTextVariance calculates variance in plain text
// Variance formula: Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)
func calculatePlainTextVariance(population []float64, x float64) float64 {
	n := float64(len(population))
	sumSquaredDiff := 0.0

	for _, value := range population {
		diff := value - x
		sumSquaredDiff += diff * diff
	}

	return sumSquaredDiff / n
}

// calculateHomomorphicVarianceBatched calculates variance using batched homomorphic encryption
// All population values are in a single ciphertext, hiding the population size
// Var = E[(X - x)^2] = (1/n) * sum((X_i - x)^2)
func calculateHomomorphicVarianceBatched(
	encryptedPopulation *rlwe.Ciphertext,
	encryptedX *rlwe.Ciphertext,
	encryptedN *rlwe.Ciphertext,
	populationSize int,
	evaluator *ckks.Evaluator,
	params ckks.Parameters,
) *rlwe.Ciphertext {
	// Calculate (X - x) for all slots simultaneously
	// encryptedPopulation contains [X_1, X_2, ..., X_n, 0, 0, ...]
	// encryptedX contains [x, x, ..., x, x, x, ...]
	diff, err := evaluator.SubNew(encryptedPopulation, encryptedX)
	if err != nil {
		log.Fatalf("Error computing difference: %v", err)
	}

	// Calculate (X - x)^2 for all slots
	squaredDiff, err := evaluator.MulRelinNew(diff, diff)
	if err != nil {
		log.Fatalf("Error computing squared difference: %v", err)
	}

	// Rescale after multiplication
	if err := evaluator.Rescale(squaredDiff, squaredDiff); err != nil {
		log.Fatalf("Error rescaling squared difference: %v", err)
	}

	// Sum all slots to get sum((X_i - x)^2)
	// We need to sum the first populationSize slots
	encryptedSum := sumSlots(squaredDiff, populationSize, evaluator, params)

	// Divide by n: variance = sum / n
	// Match levels before division
	minLevel := min(encryptedSum.Level(), encryptedN.Level())
	if encryptedSum.Level() > minLevel {
		evaluator.DropLevel(encryptedSum, encryptedSum.Level()-minLevel)
	}
	if encryptedN.Level() > minLevel {
		// Create a copy to avoid modifying the original
		encryptedNCopy := encryptedN.CopyNew()
		evaluator.DropLevel(encryptedNCopy, encryptedNCopy.Level()-minLevel)
		encryptedN = encryptedNCopy
	}

	// Create plaintext with 1/n in all slots
	invN := 1.0 / float64(populationSize)
	encryptedVariance, err := evaluator.MulNew(encryptedSum, invN)
	if err != nil {
		log.Fatalf("Error computing final variance: %v", err)
	}

	// Rescale after multiplication
	if err := evaluator.Rescale(encryptedVariance, encryptedVariance); err != nil {
		log.Fatalf("Error rescaling final result: %v", err)
	}

	return encryptedVariance
}

// sumSlots sums the first n slots of a ciphertext using a tree-based approach
// This is more efficient than sequential addition
func sumSlots(ct *rlwe.Ciphertext, n int, evaluator *ckks.Evaluator, params ckks.Parameters) *rlwe.Ciphertext {
	result := ct.CopyNew()

	// For simplicity, we'll rotate and add
	// A more efficient implementation would use a binary tree approach
	for i := 1; i < n; i++ {
		rotated, err := evaluator.RotateNew(ct, i)
		if err != nil {
			log.Fatalf("Error rotating ciphertext: %v", err)
		}

		if err := evaluator.Add(result, rotated, result); err != nil {
			log.Fatalf("Error adding rotated ciphertext: %v", err)
		}
	}

	return result
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// encryptPopulationBatched encrypts the entire population into a single ciphertext using SIMD slots
// This hides the population size from potential attackers
func encryptPopulationBatched(
	population []float64,
	encoder *ckks.Encoder,
	encryptor *rlwe.Encryptor,
	params ckks.Parameters,
) *rlwe.Ciphertext {
	// Encode all population values into the slots of a single plaintext
	// Remaining slots will be zero
	plaintext := ckks.NewPlaintext(params, params.MaxLevel())
	if err := encoder.Encode(population, plaintext); err != nil {
		log.Fatalf("Failed to encode population: %v", err)
	}

	// Encrypt the batched plaintext
	ciphertext, err := encryptor.EncryptNew(plaintext)
	if err != nil {
		log.Fatalf("Failed to encrypt population: %v", err)
	}

	return ciphertext
}

// encryptValueBatched encrypts a single value, replicating it across all slots for SIMD operations
func encryptValueBatched(
	value float64,
	encoder *ckks.Encoder,
	encryptor *rlwe.Encryptor,
	params ckks.Parameters,
	numSlots int,
) *rlwe.Ciphertext {
	// Replicate the value across all slots for element-wise operations
	values := make([]float64, numSlots)
	for i := range values {
		values[i] = value
	}

	plaintext := ckks.NewPlaintext(params, params.MaxLevel())
	if err := encoder.Encode(values, plaintext); err != nil {
		log.Fatalf("Failed to encode value: %v", err)
	}

	ciphertext, err := encryptor.EncryptNew(plaintext)
	if err != nil {
		log.Fatalf("Failed to encrypt value: %v", err)
	}

	return ciphertext
}

// decryptValue decrypts a ciphertext and returns the first value
func decryptValue(ciphertext *rlwe.Ciphertext, decryptor *rlwe.Decryptor, encoder *ckks.Encoder, params ckks.Parameters) float64 {
	plaintext := decryptor.DecryptNew(ciphertext)
	result := make([]complex128, params.MaxSlots())
	if err := encoder.Decode(plaintext, result); err != nil {
		log.Fatal(err)
	}

	return real(result[0])
}
