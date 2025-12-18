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
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         "Genesis Block",
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

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Data:         data,
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
}

// IsValid validates the integrity of the blockchain
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Validate current block's hash
		if currentBlock.Hash != currentBlock.CalculateHash() {
			fmt.Printf("Block #%d: Current hash is invalid\n", currentBlock.Index)
			return false
		}

		// Validate previous hash linking
		if currentBlock.PreviousHash != prevBlock.Hash {
			fmt.Printf("Block #%d: Previous hash is invalid\n", currentBlock.Index)
			return false
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

