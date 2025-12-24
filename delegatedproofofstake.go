package main

import (
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// Delegate represents a delegate in DPoS system
type Delegate struct {
	Address   string
	Votes     float64
	Stake     float64
	IsActive  bool
	LastBlock int
}

// DelegatedProofOfStake represents a Delegated Proof of Stake consensus mechanism
type DelegatedProofOfStake struct {
	Block     *Block
	Delegates map[string]*Delegate          // Address -> Delegate
	Votes     map[string]map[string]float64 // Voter -> Delegate -> Vote amount
}

// NewDelegatedProofOfStake creates a new DPoS instance
func NewDelegatedProofOfStake(block *Block, stakeholders map[string]float64) *DelegatedProofOfStake {
	dpos := &DelegatedProofOfStake{
		Block:     block,
		Delegates: make(map[string]*Delegate),
		Votes:     make(map[string]map[string]float64),
	}

	// Initialize delegates from stakeholders
	for address, stake := range stakeholders {
		if stake > 0 {
			dpos.Delegates[address] = &Delegate{
				Address:   address,
				Stake:     stake,
				Votes:     0,
				IsActive:  true,
				LastBlock: -1,
			}
		}
	}

	return dpos
}

// Vote allows a stakeholder to vote for a delegate
func (dpos *DelegatedProofOfStake) Vote(voterAddress string, delegateAddress string, voteAmount float64) error {
	// Check if delegate exists
	if _, exists := dpos.Delegates[delegateAddress]; !exists {
		return fmt.Errorf("delegate %s does not exist", delegateAddress)
	}

	// Initialize voter's vote map if needed
	if dpos.Votes[voterAddress] == nil {
		dpos.Votes[voterAddress] = make(map[string]float64)
	}

	// Update votes
	oldVote := dpos.Votes[voterAddress][delegateAddress]
	dpos.Delegates[delegateAddress].Votes += voteAmount - oldVote
	dpos.Votes[voterAddress][delegateAddress] = voteAmount

	return nil
}

// GetTopDelegates returns the top N delegates by votes
func (dpos *DelegatedProofOfStake) GetTopDelegates(n int) []*Delegate {
	delegates := make([]*Delegate, 0, len(dpos.Delegates))
	for _, delegate := range dpos.Delegates {
		if delegate.IsActive {
			delegates = append(delegates, delegate)
		}
	}

	// Sort by votes (descending)
	sort.Slice(delegates, func(i, j int) bool {
		return delegates[i].Votes > delegates[j].Votes
	})

	if n > len(delegates) {
		n = len(delegates)
	}

	return delegates[:n]
}

// SelectValidator selects a validator from top delegates using round-robin
func (dpos *DelegatedProofOfStake) SelectValidator() string {
	topDelegates := dpos.GetTopDelegates(21) // Top 21 delegates (common in DPoS systems)
	if len(topDelegates) == 0 {
		return ""
	}

	// Use block index for round-robin selection
	blockIndex := dpos.Block.Index
	selectedIndex := blockIndex % len(topDelegates)
	if selectedIndex < 0 {
		selectedIndex = 0
	}

	return topDelegates[selectedIndex].Address
}

// Validate validates that the validator is a valid delegate
func (dpos *DelegatedProofOfStake) Validate(validatorAddress string) bool {
	delegate, exists := dpos.Delegates[validatorAddress]
	if !exists {
		return false
	}

	if !delegate.IsActive {
		return false
	}

	// Check if delegate is in top delegates
	topDelegates := dpos.GetTopDelegates(21)
	for _, topDelegate := range topDelegates {
		if topDelegate.Address == validatorAddress {
			return true
		}
	}

	return false
}

// CalculateStakeFromVotes calculates total stake from votes
func (dpos *DelegatedProofOfStake) CalculateStakeFromVotes() map[string]float64 {
	stakes := make(map[string]float64)
	for voter, votes := range dpos.Votes {
		totalVote := 0.0
		for _, voteAmount := range votes {
			totalVote += voteAmount
		}
		stakes[voter] = totalVote
	}
	return stakes
}

// CreateBlockWithDPoS creates a block using Delegated Proof of Stake
func (bc *Blockchain) CreateBlockWithDPoS(transactions []*Transaction, validatorAddress string) error {
	// Validate all transactions
	for _, tx := range transactions {
		if err := bc.ValidateTransaction(tx); err != nil {
			return err
		}
		if tx.Signature != "" && !tx.Verify() {
			return fmt.Errorf("transaction has invalid signature")
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

	// Validate DPoS
	stakeholders := bc.CalculateStakeFromBlockchain()
	dpos := NewDelegatedProofOfStake(newBlock, stakeholders)
	if !dpos.Validate(validatorAddress) {
		return fmt.Errorf("invalid validator: %s is not a valid delegate", validatorAddress)
	}

	// Calculate hash (DPoS doesn't require mining, just hash)
	newBlock.Hash = newBlock.CalculateHash()

	bc.Blocks = append(bc.Blocks, newBlock)

	// Remove transactions from mempool
	txHashes := make([]string, 0)
	for _, tx := range transactions {
		if tx.From != "" {
			txHashes = append(txHashes, hex.EncodeToString(tx.Hash()))
		}
	}
	bc.Mempool.RemoveTransactions(txHashes)

	fmt.Printf("Block #%d added using Delegated Proof of Stake (Validator: %s)\n\n", newBlock.Index, validatorAddress[:16]+"...")
	return nil
}

// VoteForDelegate allows a stakeholder to vote for a delegate
func (bc *Blockchain) VoteForDelegate(voterAddress string, delegateAddress string, voteAmount float64) error {
	// Check if voter has sufficient balance
	balance := bc.GetTotalBalance(voterAddress)
	if balance < voteAmount {
		return fmt.Errorf("insufficient balance for voting: have %.2f, trying to vote %.2f", balance, voteAmount)
	}

	// Get current stakeholders
	stakeholders := bc.CalculateStakeFromBlockchain()
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	dpos := NewDelegatedProofOfStake(lastBlock, stakeholders)

	// Vote
	return dpos.Vote(voterAddress, delegateAddress, voteAmount)
}

// GetTopDelegates returns top delegates by votes
func (bc *Blockchain) GetTopDelegates(n int) []*Delegate {
	stakeholders := bc.CalculateStakeFromBlockchain()
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	dpos := NewDelegatedProofOfStake(lastBlock, stakeholders)

	// Initialize votes from stakes (simplified: stake = vote)
	for address, stake := range stakeholders {
		if stake > 0 {
			dpos.Vote(address, address, stake) // Self-vote with stake
		}
	}

	return dpos.GetTopDelegates(n)
}
