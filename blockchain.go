package main

import (
	"fmt"
	"time"
)

// Blockchain represents a blockchain
type Blockchain struct {
	Blocks []*Block
}

// NewBlockchain creates a new blockchain with genesis block
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks: []*Block{},
	}
	bc.CreateGenesisBlock()
	return bc
}

// CreateGenesisBlock creates the first block in the blockchain
func (bc *Blockchain) CreateGenesisBlock() {
	// Create genesis transaction
	genesisTx := NewTransaction("", "Genesis", 0)
	transactions := []*Transaction{genesisTx}

	// Create Merkle tree
	merkleTree := NewMerkleTree(transactions)
	merkleRoot := merkleTree.GetRootHash()

	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: transactions,
		MerkleRoot:   merkleRoot,
		PreviousHash: "0",
		Nonce:        0,
	}

	// Mine the genesis block
	pow := NewProofOfWork(genesisBlock)
	nonce, hash := pow.Run()
	genesisBlock.Nonce = nonce
	genesisBlock.Hash = hash

	bc.Blocks = append(bc.Blocks, genesisBlock)
	fmt.Println("Genesis block created and mined!")
}

// AddBlock adds a new block with transactions to the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) error {
	// Validate all transactions before adding
	for _, tx := range transactions {
		if err := bc.ValidateTransaction(tx); err != nil {
			return err
		}
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	// Create Merkle tree from transactions
	merkleTree := NewMerkleTree(transactions)
	merkleRoot := merkleTree.GetRootHash()

	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: transactions,
		MerkleRoot:   merkleRoot,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
	}

	// Mine the new block
	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()
	newBlock.Nonce = nonce
	newBlock.Hash = hash

	bc.Blocks = append(bc.Blocks, newBlock)
	fmt.Printf("Block #%d added to the blockchain!\n\n", newBlock.Index)
	return nil
}

// IsValid validates the integrity of the blockchain
func (bc *Blockchain) IsValid() bool {
	for i := 0; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]

		// Validate Merkle root
		merkleTree := NewMerkleTree(currentBlock.Transactions)
		calculatedMerkleRoot := merkleTree.GetRootHash()
		if currentBlock.MerkleRoot != calculatedMerkleRoot {
			fmt.Printf("Block #%d: Merkle root is invalid\n", currentBlock.Index)
			return false
		}

		// Validate transaction signatures (skip genesis block)
		if i > 0 {
			for j, tx := range currentBlock.Transactions {
				// Skip unsigned transactions (like genesis)
				if tx.Signature == "" {
					continue
				}
				// Note: In a real implementation, we'd need to store public keys
				// For now, we'll just check that signature exists and is valid format
				if len(tx.Signature) != 128 { // 64 bytes = 128 hex chars
					fmt.Printf("Block #%d: Transaction #%d has invalid signature format\n", currentBlock.Index, j+1)
					return false
				}
			}
		}

		// Validate current block's hash
		if currentBlock.Hash != currentBlock.CalculateHash() {
			fmt.Printf("Block #%d: Current hash is invalid\n", currentBlock.Index)
			return false
		}

		// Validate previous hash linking (skip genesis block)
		if i > 0 {
			prevBlock := bc.Blocks[i-1]
			if currentBlock.PreviousHash != prevBlock.Hash {
				fmt.Printf("Block #%d: Previous hash is invalid\n", currentBlock.Index)
				return false
			}
		}

		// Validate proof of work
		pow := NewProofOfWork(currentBlock)
		if !pow.Validate() {
			fmt.Printf("Block #%d: Proof of work is invalid\n", currentBlock.Index)
			return false
		}
	}

	return true
}

// Print prints all blocks in the blockchain
func (bc *Blockchain) Print() {
	for _, block := range bc.Blocks {
		fmt.Println(block.String())
		fmt.Println("---")
	}
}
