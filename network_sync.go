package main

import (
	"encoding/hex"
	"fmt"
)

// MergeBlockchain merges a received blockchain with the current one
// Uses the longest valid chain rule
func (bc *Blockchain) MergeBlockchain(receivedBlocks []*Block) error {
	if len(receivedBlocks) == 0 {
		return fmt.Errorf("received empty blockchain")
	}

	// Validate received blockchain
	if !validateBlockchain(receivedBlocks) {
		return fmt.Errorf("received blockchain is invalid")
	}

	// Use longest chain rule: if received chain is longer, replace current chain
	if len(receivedBlocks) > len(bc.Blocks) {
		bc.Blocks = receivedBlocks
		fmt.Printf("Blockchain updated: received chain is longer (%d blocks vs %d blocks)\n",
			len(receivedBlocks), len(bc.Blocks))
		return nil
	}

	// If same length, keep current chain (could add more sophisticated comparison)
	if len(receivedBlocks) == len(bc.Blocks) {
		fmt.Printf("Blockchain sync: chains have same length (%d blocks), keeping current chain\n", len(bc.Blocks))
		return nil
	}

	// Received chain is shorter, keep current chain
	fmt.Printf("Blockchain sync: current chain is longer (%d blocks vs %d blocks), keeping current chain\n",
		len(bc.Blocks), len(receivedBlocks))
	return nil
}

// validateBlockchain validates a blockchain structure
func validateBlockchain(blocks []*Block) bool {
	if len(blocks) == 0 {
		return false
	}

	// Validate genesis block
	if blocks[0].Index != 0 || blocks[0].PreviousHash != "0" {
		return false
	}

	// Validate all blocks
	for i := 0; i < len(blocks); i++ {
		currentBlock := blocks[i]

		// Validate Merkle root
		merkleTree := NewMerkleTree(currentBlock.Transactions)
		calculatedMerkleRoot := merkleTree.GetRootHash()
		if currentBlock.MerkleRoot != calculatedMerkleRoot {
			return false
		}

		// Validate transaction signatures
		if i > 0 {
			for _, tx := range currentBlock.Transactions {
				if tx.Signature != "" && !tx.Verify() {
					return false
				}
			}
		}

		// Validate hash
		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		// Validate previous hash linking
		if i > 0 {
			prevBlock := blocks[i-1]
			if currentBlock.PreviousHash != prevBlock.Hash {
				return false
			}
		}

		// Validate proof of work
		pow := NewProofOfWork(currentBlock)
		if !pow.Validate() {
			return false
		}
	}

	return true
}

// AddReceivedBlock validates and adds a block received from network
func (bc *Blockchain) AddReceivedBlock(block *Block) error {
	// Validate block structure
	if block == nil {
		return fmt.Errorf("received nil block")
	}

	// Check if block already exists
	if block.Index < len(bc.Blocks) {
		existingBlock := bc.Blocks[block.Index]
		if existingBlock.Hash == block.Hash {
			return fmt.Errorf("block already exists")
		}
	}

	// Validate block is next in sequence
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	if block.Index != prevBlock.Index+1 {
		return fmt.Errorf("block index mismatch: expected %d, got %d", prevBlock.Index+1, block.Index)
	}

	if block.PreviousHash != prevBlock.Hash {
		return fmt.Errorf("previous hash mismatch")
	}

	// Validate Merkle root
	merkleTree := NewMerkleTree(block.Transactions)
	calculatedMerkleRoot := merkleTree.GetRootHash()
	if block.MerkleRoot != calculatedMerkleRoot {
		return fmt.Errorf("invalid Merkle root")
	}

	// Validate transaction signatures
	for _, tx := range block.Transactions {
		if tx.Signature != "" && !tx.Verify() {
			return fmt.Errorf("invalid transaction signature")
		}
	}

	// Validate hash
	if block.Hash != block.CalculateHash() {
		return fmt.Errorf("invalid block hash")
	}

	// Validate proof of work
	pow := NewProofOfWork(block)
	if !pow.Validate() {
		return fmt.Errorf("invalid proof of work")
	}

	// Add block to blockchain
	bc.Blocks = append(bc.Blocks, block)

	// Remove transactions from mempool
	txHashes := make([]string, 0)
	for _, tx := range block.Transactions {
		// Skip reward transactions
		if tx.From == "" {
			continue
		}
		txHashes = append(txHashes, hex.EncodeToString(tx.Hash()))
	}
	bc.Mempool.RemoveTransactions(txHashes)

	fmt.Printf("Block #%d added from network\n", block.Index)
	return nil
}
