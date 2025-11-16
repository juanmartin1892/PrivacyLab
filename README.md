# PrivacyLab

## Purpose

**PrivacyLab** is a hands-on lab designed to explore, prototype and demonstrate **advanced secure computing techniques**, combining:

- Homomorphic encryption
- Zero-knowledge proofs
- Isolation and separation of roles
- Secure key and identity architectures
- Multiparty computing in controlled contexts
- Verifiable end-to-end integrity
- Good engineering practices (SOLID + Hexagonal Architecture)

The goal is not to build a commercial product, but to create a controlled environment in which to systematically experiment and learn about these mechanisms.

## Vision

PrivacyLab's vision is to serve as a conceptual and practical basis for systems where:

- A **client** can execute operations on its data without disclosing it to the server.
- The **integrity** and **privacy** of data and results are demonstrably guaranteed.
- Results are **cryptographically verifiable**, without relying solely on trust in the infrastructure or the operator.
- The architecture is **modular, extensible and testable**, ready to integrate new secure computing techniques as they are explored.

## Current Scope

This project is in its **initial phase**. The current focus is on:

- Establishing the foundational architecture and project structure.
- Defining clear motivations and guiding principles.
- Setting up the development environment and tooling.
- Implementing initial prototypes and experiments with secure computing techniques.

**Status:** Active development with working examples.

### Implemented Features

#### Point-to-Point Variance Calculation with Homomorphic Encryption

A complete working example demonstrating privacy-preserving variance calculation using the CKKS homomorphic encryption scheme.

**Key Features:**
- **SIMD Batching**: All population values encrypted in a single ciphertext for enhanced privacy.
- **Hidden Population Size**: Prevents attackers from determining the number of individuals by counting ciphertexts.
- **Fully Homomorphic Operations**: Variance computed entirely on encrypted data without decryption.
- **Verifiable Results**: Plain-text comparison ensures computational correctness.

**Technical Highlights:**
- Uses Lattigo v6 library for homomorphic encryption.
- Implements CKKS scheme for approximate arithmetic on encrypted real numbers.
- Demonstrates rotation keys for SIMD operations.
- Achieves high precision with minimal error (relative error < 0.000001%).

**Location:** `examples/point-to-point/`

**Quick Start:**
```bash
cd examples/point-to-point
./run.sh
```

See `examples/point-to-point/README.md` for detailed documentation and execution instructions.

## Project Structure

```
privacyLab/
├── cmd/                    # Application entry points
│   ├── client/            # Client application
│   └── worker/            # Worker application
├── examples/              # Working examples and demonstrations
│   └── point-to-point/   # Variance calculation with homomorphic encryption
├── internal/              # Internal packages
│   ├── adapters/         # External interfaces and implementations
│   ├── domain/           # Core business logic and domain models
│   └── ports/            # Interface definitions
├── go.mod                # Go module dependencies
└── README.md             # This file
```

## Examples and Demonstrations

### 1. Point-to-Point Variance Calculation

Demonstrates privacy-preserving statistical computation using homomorphic encryption.

- **Location:** `examples/point-to-point/`
- **Technology:** CKKS homomorphic encryption (Lattigo v6)
- **Use Case:** Calculate variance of an input value against an encrypted population
- **Privacy Features:** SIMD batching, hidden population size, fully encrypted computation
- **Documentation:** See `examples/point-to-point/README.md`

**Run it:**
```bash
cd examples/point-to-point && ./run.sh
```

## Motivations

- To explore in a practical way how to combine homomorphic encryption, zero-knowledge proofs and MPC to **minimize information disclosure**.
- To understand the **performance and complexity limits** of these techniques in a controlled environment.
- Design and document **architectural patterns** to serve as a reference for future sensitive data computing systems.
- Build a reusable knowledge base, with clear documentation, that facilitates the introduction of new use cases or new cryptographic primitives.

## Author

**Author:** Juan Martín Pérez

Personal research and exploration project in secure computing.