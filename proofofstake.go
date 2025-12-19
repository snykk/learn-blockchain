package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"
)

// ProofOfStake represents a proof of stake consensus mechanism
type ProofOfStake struct {
	Block        *Block
	Stakeholders map[string]float64 // Address -> Stake amount
}

// NewProofOfStake creates a new proof of stake
func NewProofOfStake(block *Block, stakeholders map[string]float64) *ProofOfStake {
	return &ProofOfStake{
		Block:        block,
		Stakeholders: stakeholders,
	}
}

// SelectValidator selects a validator based on stake (weighted random selection)
func (pos *ProofOfStake) SelectValidator() string {
	if len(pos.Stakeholders) == 0 {
		return ""
	}

	// Calculate total stake
	totalStake := 0.0
	for _, stake := range pos.Stakeholders {
		totalStake += stake
	}

	if totalStake == 0 {
		return ""
	}

	// Use block hash as seed for deterministic selection
	seed := pos.Block.PreviousHash + pos.Block.MerkleRoot
	hash := sha256.Sum256([]byte(seed))
	hashInt := new(big.Int).SetBytes(hash[:])

	// Convert to float64 for percentage calculation
	// Use modulo to get a value between 0 and totalStake
	modValue := new(big.Float).SetInt(hashInt)
	modValue.Quo(modValue, big.NewFloat(1e38)) // Normalize
	randomValue, _ := modValue.Float64()
	randomValue = randomValue * totalStake

	// Select validator based on weighted stake
	currentSum := 0.0
	for address, stake := range pos.Stakeholders {
		currentSum += stake
		if randomValue <= currentSum {
			return address
		}
	}

	// Fallback: return first stakeholder
	for address := range pos.Stakeholders {
		return address
	}

	return ""
}

// Validate validates the proof of stake
func (pos *ProofOfStake) Validate(validatorAddress string) bool {
	// Check if validator has stake
	stake, exists := pos.Stakeholders[validatorAddress]
	if !exists || stake <= 0 {
		return false
	}

	// Validate that the validator was selected correctly
	selectedValidator := pos.SelectValidator()
	return selectedValidator == validatorAddress
}

// CalculateStakeFromBlockchain calculates stakeholder stakes from blockchain balances
func (bc *Blockchain) CalculateStakeFromBlockchain() map[string]float64 {
	stakeholders := make(map[string]float64)

	// Get all unique addresses from transactions
	addresses := make(map[string]bool)
	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.From != "" {
				addresses[tx.From] = true
			}
			if tx.To != "" && tx.To != "Genesis" {
				addresses[tx.To] = true
			}
		}
	}

	// Calculate stake as balance
	for address := range addresses {
		balance := bc.GetBalance(address)
		if balance > 0 {
			stakeholders[address] = balance
		}
	}

	return stakeholders
}

// CreateBlockWithPoS creates a block using Proof of Stake instead of Proof of Work
func (bc *Blockchain) CreateBlockWithPoS(transactions []*Transaction, validatorAddress string) error {
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
		Nonce:        0, // PoS doesn't use nonce for mining
	}

	// Validate Proof of Stake
	stakeholders := bc.CalculateStakeFromBlockchain()
	pos := NewProofOfStake(newBlock, stakeholders)
	if !pos.Validate(validatorAddress) {
		return fmt.Errorf("invalid validator: %s does not have sufficient stake or was not selected", validatorAddress)
	}

	// Calculate hash (PoS doesn't require mining, just hash)
	newBlock.Hash = newBlock.CalculateHash()

	bc.Blocks = append(bc.Blocks, newBlock)
	fmt.Printf("Block #%d added to the blockchain using Proof of Stake! (Validator: %s)\n\n", newBlock.Index, validatorAddress)
	return nil
}
