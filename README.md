# Learn Blockchain

An enhanced blockchain implementation in Go that covers fundamental blockchain concepts including transactions, Merkle trees, and wallet signing.

## Implemented Concepts

1. **Block Structure** - Block data structure with index, timestamp, transactions, Merkle root, previous hash, hash, and nonce
2. **Cryptographic Hashing** - SHA-256 to generate unique hash for each block
3. **Chain Linking** - Each block is linked to the previous block through previous hash
4. **Proof of Work (Mining)** - Consensus mechanism with difficulty target
5. **Genesis Block** - The first block in the chain
6. **Blockchain Validation** - Chain integrity validation (hash validation, chain linking, proof of work, Merkle root)
7. **Transactions** - Transaction structure with sender, receiver, amount, and digital signature
8. **Merkle Tree** - Efficient transaction verification using Merkle tree structure
9. **Wallet & Signing** - ECDSA-based wallet with transaction signing and verification

## File Structure

```
learn-blockchain/
├── README.md           # Documentation
├── go.mod              # Go module file
├── main.go             # Entry point and demo
├── block.go            # Block structure and methods
├── blockchain.go       # Blockchain structure and methods
├── proofofwork.go      # Proof of Work implementation
├── transaction.go      # Transaction structure and signing
├── merkle.go           # Merkle tree implementation
├── wallet.go           # Wallet with ECDSA key generation
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
- **Transactions**: Array of transactions stored in the block
- **MerkleRoot**: Root hash of the Merkle tree built from transactions
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
- Contains a genesis transaction
- Must be mined like other blocks

### 6. Blockchain Validation

Blockchain validation checks:
- **Merkle Root Validation**: Whether each block's Merkle root matches the calculated root from transactions
- **Transaction Signature Validation**: Whether transaction signatures are valid (format check)
- **Hash Validation**: Whether each block's hash matches the calculated hash from block data
- **Chain Linking**: Whether each block's `PreviousHash` matches the previous block's hash
- **Proof of Work**: Whether each block meets the target difficulty

### 7. Transactions

Transactions represent value transfers in the blockchain:
- **From**: Sender's wallet address
- **To**: Receiver's wallet address
- **Amount**: Amount being transferred
- **Signature**: Digital signature created using sender's private key (ECDSA)

### 8. Merkle Tree

Merkle tree provides efficient transaction verification:
- All transactions in a block are hashed and organized in a binary tree
- Root hash (Merkle root) is stored in the block header
- Allows efficient verification of transaction inclusion without downloading all transactions
- Any change in transactions will result in a different Merkle root

### 9. Wallet & Signing

Wallets provide cryptographic key management:
- **Key Generation**: ECDSA key pair generation using P-256 curve
- **Address Generation**: Wallet address derived from public key hash
- **Transaction Signing**: Transactions are signed with private key before being added to blocks
- **Signature Verification**: Transaction signatures can be verified using public key

## Example Output

The program will display:
1. Wallet creation for multiple users (Alice, Bob, Charlie)
2. Blockchain creation with genesis block
3. Creating and signing transactions
4. Adding blocks with signed transactions
5. Displaying the entire blockchain with transactions
6. Blockchain validation (including Merkle root and signatures)
7. Tampering detection tests (modify transaction and check validation)

## Features

### Core Blockchain Features
- Block structure with all important fields
- SHA-256 hashing
- Proof of Work with adjustable difficulty (targetBits)
- Genesis block creation
- Chain linking (previous hash)
- Blockchain validation
- Tampering detection
- CLI demo for testing

### Enhanced Features
- **Transaction System**: Structured transactions with sender, receiver, and amount
- **Merkle Tree**: Efficient transaction verification using Merkle tree structure
- **Wallet System**: ECDSA-based wallet with key generation and address creation
- **Digital Signatures**: Transaction signing and verification using ECDSA cryptography
- **Transaction Validation**: Merkle root validation and signature format checking

## Adjusting Difficulty

To change mining difficulty, edit the `targetBits` constant in `proofofwork.go`:

```go
const targetBits = 16 // Difficulty: hash must start with 4 leading zeros
```

- Smaller value = easier (faster mining)
- Larger value = harder (slower mining)

## Enhanced Features Implemented

This implementation now includes:
- **Merkle Tree** for transactions
- **Wallet and transaction signing** with ECDSA
- **Transaction structure** with digital signatures

## Future Enhancements

- Network/P2P for distributed blockchain
- Smart Contracts
- Different consensus mechanisms (Proof of Stake, etc.)
- Full signature verification with public key storage
- Transaction pool/mempool
- Web3 integration
