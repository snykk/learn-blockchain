# Learn Blockchain

A simple blockchain implementation in Go that covers all fundamental blockchain concepts.

## Implemented Concepts

1. **Block Structure** - Block data structure with index, timestamp, data, previous hash, hash, and nonce
2. **Cryptographic Hashing** - SHA-256 to generate unique hash for each block
3. **Chain Linking** - Each block is linked to the previous block through previous hash
4. **Proof of Work (Mining)** - Consensus mechanism with difficulty target
5. **Genesis Block** - The first block in the chain
6. **Blockchain Validation** - Chain integrity validation (hash validation, chain linking, proof of work)

## File Structure

```
learn-blockchain/
├── README.md           # Documentation
├── go.mod              # Go module file
├── main.go             # Entry point and demo
├── block.go            # Block structure and methods
├── blockchain.go       # Blockchain structure and methods
├── proofofwork.go      # Proof of Work implementation
└── utils.go            # Utility functions (hashing, etc.)
```

## How to Run

### Prerequisites

- Go 1.16 or newer

### Running the Program

```bash
go run .
```

Or build first:

```bash
go build
./learn-blockchain
```

## Concept Explanation

### 1. Block Structure

Each block in the blockchain has the following structure:

- **Index**: Position of the block in the chain (starts from 0 for genesis block)
- **Timestamp**: Time when the block was created
- **Data**: Data or transactions stored in the block
- **PreviousHash**: Hash of the previous block (links blocks in the chain)
- **Hash**: Hash of this block (calculated from all fields including nonce)
- **Nonce**: Number used once, value used in mining to find a valid hash

### 2. Cryptographic Hashing

Uses SHA-256 to generate a unique hash for each block. This hash is:
- Deterministic: same input always produces the same hash
- One-way: cannot be reversed to get the original data
- Avalanche effect: small changes in input produce vastly different hashes

### 3. Chain Linking

Each block stores the hash of the previous block in the `PreviousHash` field. This creates an immutable chain:
- If data in any block is changed, its hash will change
- Hash change will make the `PreviousHash` in the next block invalid
- The entire chain after the modified block becomes invalid

### 4. Proof of Work (Mining)

Proof of Work is a consensus mechanism that:
- Requires computation to find a nonce that produces a hash with specific characteristics
- Difficulty is determined by the number of leading zeros required in the hash
- Higher difficulty means longer time required for mining
- Prevents spam and ensures blockchain security

### 5. Genesis Block

The genesis block is the first block in the blockchain:
- Index = 0
- PreviousHash = "0" (because there's no previous block)
- Data = "Genesis Block"
- Must be mined like other blocks

### 6. Blockchain Validation

Blockchain validation checks:
- **Hash Validation**: Whether each block's hash matches the calculated hash from block data
- **Chain Linking**: Whether each block's `PreviousHash` matches the previous block's hash
- **Proof of Work**: Whether each block meets the target difficulty

## Example Output

The program will display:
1. Blockchain creation with genesis block
2. Adding several blocks with transaction data
3. Displaying the entire blockchain
4. Blockchain validation
5. Tampering detection test (modify data and check validation)

## Features

- Block structure with all important fields
- SHA-256 hashing
- Proof of Work with adjustable difficulty (targetBits)
- Genesis block creation
- Chain linking (previous hash)
- Blockchain validation
- Tampering detection
- CLI demo for testing

## Adjusting Difficulty

To change mining difficulty, edit the `targetBits` constant in `proofofwork.go`:

```go
const targetBits = 16 // Difficulty: hash must start with 4 leading zeros
```

- Smaller value = easier (faster mining)
- Larger value = harder (slower mining)

## Learn More

This implementation covers basic blockchain concepts. For further development, consider:

- Merkle Tree for transactions
- Network/P2P for distributed blockchain
- Smart Contracts
- Different consensus mechanisms (Proof of Stake, etc.)
- Wallet and transaction signing
- Web3 integration
