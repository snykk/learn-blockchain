# Learn Blockchain

An enhanced blockchain implementation in Go that covers fundamental blockchain concepts including transactions, Merkle trees, wallet signing, mempool, full signature verification, and Proof of Stake consensus.

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
10. **Transaction Pool (Mempool)** - Pending transaction storage before block creation
11. **Full Signature Verification** - Complete signature verification using stored public keys
12. **Proof of Stake** - Alternative consensus mechanism based on stake weight
13. **Delegated Proof of Stake (DPoS)** - Consensus mechanism with delegate voting and round-robin selection
14. **Balance System** - Balance tracking and validation for transactions
15. **Transaction Fees** - Fees paid by transaction senders for processing
16. **Block Rewards** - Rewards given to miners/validators for creating blocks
17. **Network/P2P** - Peer-to-peer network for distributed blockchain
18. **Smart Contracts** - Executable contracts deployed on the blockchain
19. **Web3 Integration** - Web3 JSON-RPC API server with Ethereum-compatible endpoints
20. **PBFT Consensus** - Practical Byzantine Fault Tolerance consensus mechanism
21. **Raft Consensus** - Leader-based consensus algorithm for distributed systems

## File Structure

```
learn-blockchain/
├── README.md           # Documentation
├── go.mod              # Go module file
├── main.go             # Entry point and demo
├── block.go            # Block structure and methods
├── blockchain.go       # Blockchain structure and methods
├── proofofwork.go      # Proof of Work implementation
├── proofofstake.go     # Proof of Stake implementation
├── delegatedproofofstake.go # Delegated Proof of Stake implementation
├── pbft.go             # PBFT consensus implementation
├── raft.go             # Raft consensus implementation
├── transaction.go      # Transaction structure and signing
├── merkle.go           # Merkle tree implementation
├── wallet.go           # Wallet with ECDSA key generation
├── mempool.go          # Transaction pool/mempool implementation
├── balance.go          # Balance calculation and validation
├── rewards.go          # Block rewards and miner rewards calculation
├── smartcontract.go    # Smart contract implementation
├── web3.go             # Web3 JSON-RPC server
├── network.go          # P2P network implementation
├── network_sync.go     # Network synchronization
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

### 10. Transaction Pool (Mempool)

Mempool stores pending transactions before they are added to blocks:
- **Transaction Storage**: Pending transactions are stored in memory pool
- **Transaction Management**: Add, remove, and retrieve transactions from mempool
- **Block Creation**: Blocks can be created from transactions in mempool
- **Automatic Cleanup**: Transactions are automatically removed from mempool when added to blocks

### 11. Full Signature Verification

Complete signature verification system:
- **Public Key Storage**: Public keys are stored in transactions for verification
- **Automatic Verification**: Signatures are verified automatically during block validation
- **ECDSA Verification**: Full ECDSA signature verification using stored public keys
- **Security**: Ensures transaction authenticity and prevents tampering

### 12. Proof of Stake

Alternative consensus mechanism to Proof of Work:
- **Stake-Based Selection**: Validators are selected based on their stake (balance)
- **Weighted Random**: Validator selection uses weighted random based on stake amount
- **No Mining Required**: PoS doesn't require computational mining like PoW
- **Energy Efficient**: More energy-efficient than Proof of Work
- **Stake Calculation**: Stake is calculated from blockchain balances

### 13. Delegated Proof of Stake (DPoS)

Advanced consensus mechanism with delegate voting:
- **Delegate System**: Stakeholders vote for delegates who validate blocks
- **Top Delegates**: System selects top N delegates (e.g., top 21) based on votes
- **Round-Robin Selection**: Validators are selected in round-robin fashion from top delegates
- **Voting Mechanism**: Stakeholders can vote for delegates using their stake
- **Efficient**: More scalable than PoS with faster block times
- **Democratic**: Stakeholders have voting power proportional to their stake

### 14. Balance System

Balance tracking and validation:
- **Balance Calculation**: Calculates balance by scanning all transactions
- **Transaction Validation**: Validates transactions before adding to blocks
- **Insufficient Balance Detection**: Prevents transactions with insufficient balance
- **Coinbase Support**: Supports coinbase transactions for initial balance distribution

### 15. Transaction Fees

Transaction fees incentivize miners and prevent spam:
- **Fee Field**: Each transaction can include a fee paid by the sender
- **Fee Deduction**: Fees are deducted from sender's balance along with transaction amount
- **Total Cost**: Total cost = Amount + Fee
- **Fee Collection**: Fees are collected by miners who create blocks
- **Optional Fees**: Transactions can be created with or without fees

### 16. Block Rewards

Block rewards incentivize miners/validators to secure the network:
- **Block Reward**: Fixed reward (50 coins) given to miner for each block created
- **Genesis Reward**: Special reward (100 coins) for genesis block
- **Reward Transaction**: Block reward is added as first transaction in block
- **Miner Rewards**: Total rewards = Block rewards + Transaction fees collected
- **Incentive Mechanism**: Rewards encourage participation in network security

### 17. Network/P2P

Peer-to-peer network for distributed blockchain:
- **Node Structure**: Each node has its own blockchain and peer list
- **TCP Communication**: Nodes communicate via TCP connections
- **Message Protocol**: JSON-based message protocol for data exchange
- **Peer Management**: Add/remove peers, maintain peer connections
- **Blockchain Synchronization**: Broadcast and sync blockchain with peers
- **Message Types**: Support for blockchain, block, transaction, ping/pong messages
- **Network Server**: Each node can run as a server to accept connections

### 18. Smart Contracts

Executable contracts deployed on the blockchain:
- **Contract Types**: Support for Simple Storage, Token (ERC-20 like), Escrow, and Voting contracts
- **Contract Deployment**: Deploy contracts with unique addresses derived from deployer and block index
- **Contract Execution**: Execute contract functions with arguments and value
- **Contract State**: Persistent state storage for each contract
- **Contract Registry**: Central registry to manage all deployed contracts
- **Contract Calls**: Transactions can include contract call data to interact with smart contracts
- **State Management**: Each contract maintains its own state that persists across calls

### 19. Web3 Integration

Web3 JSON-RPC API server for Ethereum-compatible blockchain interaction:
- **HTTP JSON-RPC Server**: Standard Web3 API server on port 8545
- **Ethereum-Compatible Endpoints**: Implements standard Ethereum JSON-RPC methods
- **API Methods**:
  - `web3_clientVersion` - Returns client version information
  - `eth_blockNumber` - Returns the latest block number
  - `eth_getBalance` - Gets account balance in Wei
  - `eth_getBlockByNumber` - Retrieves block by number
  - `eth_getTransactionCount` - Gets transaction count for address
  - `eth_sendTransaction` - Sends new transaction to mempool
  - `eth_call` - Executes contract call (read-only)
  - `eth_getCode` - Gets contract bytecode
- **Web3 Compatibility**: Compatible with Web3 libraries and tools
- **JSON-RPC 2.0**: Follows JSON-RPC 2.0 specification

### 20. PBFT Consensus

Practical Byzantine Fault Tolerance consensus mechanism:
- **Three-Phase Protocol**: Pre-Prepare, Prepare, and Commit phases
- **Byzantine Fault Tolerance**: Tolerates up to f faulty nodes in 3f+1 network
- **Quorum-Based**: Requires 2f+1 votes for consensus
- **Deterministic Finality**: No forks, immediate block finality
- **Primary Node Selection**: Round-robin primary selection based on view
- **Message Types**: Pre-Prepare, Prepare, Commit, and View Change messages
- **State Machine**: Idle → Pre-Prepare → Prepare → Commit → Finalized
- **Consensus Validation**: Validates quorum requirements and message signatures

### 21. Raft Consensus

Leader-based consensus algorithm for distributed systems:
- **Leader Election**: Democratic leader selection through voting
- **Log Replication**: Strong consistency through log replication
- **Safety Properties**: At most one leader per term, election safety
- **Heartbeat Mechanism**: Leader maintains authority through periodic heartbeats
- **Crash Fault Tolerance**: Tolerates up to (N-1)/2 node failures
- **State Transitions**: Follower → Candidate → Leader
- **Message Types**: RequestVote, RequestVoteResponse, AppendEntries, AppendEntriesResponse
- **Consensus Validation**: Majority-based replication and commit

#### Raft vs PBFT Comparison:
- **Raft**: Simpler to understand and implement, leader-based, crash fault tolerant
- **PBFT**: More complex, no fixed leader, Byzantine fault tolerant (handles malicious nodes)
- **Raft**: Better for crash failures in trusted networks
- **PBFT**: Better for malicious nodes in untrusted networks
- **Raft**: Uses randomized election timeouts (150-300ms)
- **PBFT**: Uses deterministic view-based primary selection
- **Raft**: Log replication with consistency check
- **PBFT**: Three-phase protocol with voting

#### Raft Key Components:
1. **Leader Election**:
   - Nodes start as followers
   - If no heartbeat received within election timeout, become candidate
   - Candidate requests votes from all nodes
   - Node with majority votes becomes leader

2. **Log Replication**:
   - Leader accepts transactions and creates log entries
   - Leader replicates log entries to all followers
   - Once majority acknowledges entry, it's committed
   - Leader notifies followers of commit index

3. **Safety**:
   - Election safety: at most one leader per term
   - Leader append-only: leader never overwrites/deletes log entries
   - Log matching: if two logs contain same entry at same index, all previous entries match
   - Leader completeness: if log entry is committed in one term, it appears in all future leaders' logs

4. **Heartbeat**:
   - Leader sends periodic heartbeats (empty AppendEntries) to followers
   - Followers reset election timeout on receiving heartbeat
   - Prevents unnecessary elections when leader is alive

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
- **Transaction Pool (Mempool)**: Pending transaction storage and management
- **Full Signature Verification**: Complete signature verification with stored public keys
- **Proof of Stake**: Alternative consensus mechanism based on stake weight
- **Delegated Proof of Stake**: Advanced consensus with delegate voting and round-robin selection
- **Balance System**: Balance tracking, validation, and coinbase support
- **Transaction Fees**: Optional fees for transaction processing
- **Block Rewards**: Rewards for miners/validators who create blocks
- **Network/P2P**: Peer-to-peer network for distributed blockchain
- **Smart Contracts**: Executable contracts with state management and function calls
- **Web3 Integration**: JSON-RPC API server with Ethereum-compatible endpoints
- **PBFT Consensus**: Byzantine fault-tolerant consensus mechanism

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
- **Transaction Pool (Mempool)** for pending transaction management
- **Full Signature Verification** with public key storage
- **Proof of Stake** as alternative consensus mechanism
- **Delegated Proof of Stake** with delegate voting system
- **Balance System** with validation and coinbase support
- **Transaction Fees** for transaction processing
- **Block Rewards** for miners/validators
- **Network/P2P** for distributed blockchain
- **Smart Contracts** with multiple contract types and state management
- **Web3 Integration** with JSON-RPC API server
- **PBFT Consensus** for Byzantine fault tolerance
- **Raft Consensus** for leader-based distributed consensus

## Future Enhancements

- Additional consensus mechanisms (HoneyBadgerBFT, etc.)
- Layer 2 scaling solutions
- Cross-chain bridges
- Enhanced smart contract language
