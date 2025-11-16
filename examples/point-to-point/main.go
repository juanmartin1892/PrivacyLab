package main

import (
	"fmt"
	"log"
	"math"

	"github.com/juanmartin/privacylab/internal/adapters/file"
	"github.com/juanmartin/privacylab/internal/adapters/lattigo"
	"github.com/juanmartin/privacylab/internal/domain/crypto"
	"github.com/juanmartin/privacylab/internal/domain/dataset"
	"github.com/juanmartin/privacylab/internal/domain/operation"
	"github.com/juanmartin/privacylab/internal/services/decrypt"
	"github.com/juanmartin/privacylab/internal/services/encrypt"
	"github.com/juanmartin/privacylab/internal/services/process"
)

func main() {
	fmt.Println("=== Variance Calculation with Homomorphic Encryption ===")
	fmt.Println("Computing variance of input x against population X using encrypted data")
	fmt.Println()

	// Read population from file using the file adapter
	fmt.Println("Reading population from file...")
	fileReader := file.NewFileReaderAdapter()
	populationDataset, err := fileReader.ReadFromFile("population.txt")
	if err != nil {
		log.Fatalf("Error reading population file: %v", err)
	}
	population := populationDataset.Values()
	fmt.Printf("Population loaded: %v (size: %d)\n", population, len(population))

	// Input value x
	var inputX float64
	fmt.Print("\nEnter input value x: ")
	_, err = fmt.Scanf("%f", &inputX)
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	fmt.Printf("Input x: %.2f\n", inputX)

	// Configure CKKS parameters using domain model
	fmt.Println("\nConfiguring cryptographic parameters...")
	params, err := crypto.NewParameters(
		crypto.SchemeCKKS,
		14,
		[]int{55, 40, 40, 40, 40, 40},
		[]int{61},
		40,
		128,
	)
	if err != nil {
		log.Fatalf("Error creating parameters: %v", err)
	}

	// Create variance operation to determine required rotations
	varianceOp := operation.NewVarianceOperation()
	rotations := varianceOp.RequiredRotations(populationDataset.Size())

	// Generate keys using the key generator adapter
	fmt.Println("Generating cryptographic keys...")
	keyGen, err := lattigo.NewKeyGeneratorAdapter(params)
	if err != nil {
		log.Fatalf("Error creating key generator: %v", err)
	}

	keySet, err := keyGen.GenerateKeys(params, rotations)
	if err != nil {
		log.Fatalf("Error generating keys: %v", err)
	}
	fmt.Println("Keys generated successfully")

	// Create encoder, encryptor, decryptor adapters
	fmt.Println("Setting up cryptographic services...")
	encoder, err := lattigo.NewEncoderAdapter(params)
	if err != nil {
		log.Fatalf("Error creating encoder: %v", err)
	}

	encryptor, err := lattigo.NewEncryptorAdapter(params, keySet.PublicKey())
	if err != nil {
		log.Fatalf("Error creating encryptor: %v", err)
	}

	decryptor, err := lattigo.NewDecryptorAdapter(params, keySet.SecretKey())
	if err != nil {
		log.Fatalf("Error creating decryptor: %v", err)
	}

	// Create encryption service
	encryptService := encrypt.NewService(encoder, encryptor, keySet, params)

	// Calculate variance in plain text for verification using decrypt service
	fmt.Println("\nCalculating variance in plain text (for verification)...")
	decryptService := decrypt.NewService(encoder, decryptor)
	operationParams := operation.NewOperationParams(inputX)
	plainTextVariance, err := decryptService.ComputePlaintext(varianceOp, population, operationParams)
	if err != nil {
		log.Fatalf("Error computing plaintext variance: %v", err)
	}
	fmt.Printf("Plain text variance: %.6f\n", plainTextVariance)

	// Prepare computation request using encryption service
	fmt.Println("\nEncrypting data and preparing computation request...")
	request, err := encryptService.PrepareRequest(
		operation.OperationVariance,
		[]*dataset.Dataset{populationDataset},
		[]float64{inputX, float64(populationDataset.Size())}, // Include n as third encrypted input
		map[string]interface{}{
			"data_size": populationDataset.Size(),
		},
		"variance-example-1",
	)
	if err != nil {
		log.Fatalf("Error preparing request: %v", err)
	}
	fmt.Printf("Request prepared with %d encrypted inputs\n", len(request.EncryptedInputs()))

	// Create evaluator and process service (worker side)
	fmt.Println("\nSetting up worker services...")
	evaluator, err := lattigo.NewEvaluatorAdapter(params, keySet, "worker-1")
	if err != nil {
		log.Fatalf("Error creating evaluator: %v", err)
	}

	processService := process.NewService(evaluator, "worker-1")

	// Process the encrypted request
	fmt.Println("Processing encrypted computation request...")
	result, err := processService.ProcessRequest(request)
	if err != nil {
		log.Fatalf("Error processing request: %v", err)
	}
	fmt.Printf("Computation completed in %v\n", result.ComputationTime())

	// Decrypt result using decryption service
	fmt.Println("\nDecrypting result...")
	homomorphicVariance, err := decryptService.DecryptResult(result)
	if err != nil {
		log.Fatalf("Error decrypting result: %v", err)
	}

	// Display results
	fmt.Printf("\nResults:\n")
	fmt.Printf("Plain text variance:      %.6f\n", plainTextVariance)
	fmt.Printf("Homomorphic variance:     %.6f\n", homomorphicVariance)
	fmt.Printf("Absolute error:           %.9f\n", math.Abs(homomorphicVariance-plainTextVariance))
	fmt.Printf("Relative error:           %.6f%%\n", math.Abs(homomorphicVariance-plainTextVariance)/plainTextVariance*100)

	// Display computation metadata
	fmt.Printf("\nComputation Metadata:\n")
	fmt.Printf("Worker ID:                %s\n", result.Metadata().WorkerID())
	fmt.Printf("Multiplications used:     %d\n", result.Metadata().MultiplicationsUsed())
	fmt.Printf("Rotations used:           %d\n", result.Metadata().RotationsUsed())
	fmt.Printf("Final level:              %d\n", result.Metadata().FinalLevel())
	fmt.Printf("Final scale:              %.2f\n", result.Metadata().FinalScale())
	fmt.Printf("Computation time:         %v\n", result.ComputationTime())

	fmt.Println("\nExample completed successfully")
	fmt.Println("Variance was calculated without revealing the original data")
	fmt.Println("All cryptographic operations were performed through internal services")
}
