package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// PBFTMessageType represents the type of PBFT message
type PBFTMessageType string

const (
	PrePrepare PBFTMessageType = "pre-prepare"
	Prepare    PBFTMessageType = "prepare"
	Commit     PBFTMessageType = "commit"
	ViewChange PBFTMessageType = "view-change"
)

// PBFTMessage represents a PBFT consensus message
type PBFTMessage struct {
	Type      PBFTMessageType `json:"type"`
	BlockHash string          `json:"block_hash"`
	NodeID    string          `json:"node_id"`
	Sequence  int64           `json:"sequence"`
	ViewID    int64           `json:"view_id"`
	Timestamp time.Time       `json:"timestamp"`
	Signature string          `json:"signature"`
}

// PBFTState represents the current state of PBFT consensus
type PBFTState string

const (
	StateIdle       PBFTState = "idle"
	StatePrePrepare PBFTState = "pre-prepare"
	StatePrepare    PBFTState = "prepare"
	StateCommit     PBFTState = "commit"
	StateFinalized  PBFTState = "finalized"
)

// PBFT represents a PBFT consensus instance
type PBFT struct {
	NodeID        string
	Nodes         []string // All nodes in the network
	Block         *Block
	State         PBFTState
	ViewID        int64
	Sequence      int64
	Messages      []*PBFTMessage
	PrepareCount  int
	CommitCount   int
	RequiredVotes int // 2f+1 where f is max faulty nodes
	TotalNodes    int // 3f+1
	mu            sync.RWMutex
	PrePrepared   bool
	Prepared      bool
	Committed     bool
}

// NewPBFT creates a new PBFT instance
func NewPBFT(nodeID string, nodes []string, block *Block, sequence int64) *PBFT {
	totalNodes := len(nodes)
	// In PBFT, we need 3f+1 nodes where f is the number of faulty nodes
	// RequiredVotes = 2f+1 (quorum)
	f := (totalNodes - 1) / 3
	requiredVotes := 2*f + 1

	return &PBFT{
		NodeID:        nodeID,
		Nodes:         nodes,
		Block:         block,
		State:         StateIdle,
		ViewID:        0,
		Sequence:      sequence,
		Messages:      make([]*PBFTMessage, 0),
		PrepareCount:  0,
		CommitCount:   0,
		RequiredVotes: requiredVotes,
		TotalNodes:    totalNodes,
		PrePrepared:   false,
		Prepared:      false,
		Committed:     false,
	}
}

// GetPrimaryNode returns the primary node for the current view
func (pbft *PBFT) GetPrimaryNode() string {
	primaryIndex := int(pbft.ViewID) % len(pbft.Nodes)
	return pbft.Nodes[primaryIndex]
}

// IsPrimary checks if this node is the primary for current view
func (pbft *PBFT) IsPrimary() bool {
	return pbft.NodeID == pbft.GetPrimaryNode()
}

// PrePreparePhase initiates the pre-prepare phase (only by primary)
func (pbft *PBFT) PrePreparePhase() (*PBFTMessage, error) {
	pbft.mu.Lock()
	defer pbft.mu.Unlock()

	if !pbft.IsPrimary() {
		return nil, fmt.Errorf("only primary node can initiate pre-prepare")
	}

	if pbft.State != StateIdle {
		return nil, fmt.Errorf("invalid state for pre-prepare")
	}

	msg := &PBFTMessage{
		Type:      PrePrepare,
		BlockHash: pbft.Block.Hash,
		NodeID:    pbft.NodeID,
		Sequence:  pbft.Sequence,
		ViewID:    pbft.ViewID,
		Timestamp: time.Now(),
		Signature: pbft.signMessage(PrePrepare, pbft.Block.Hash),
	}

	pbft.Messages = append(pbft.Messages, msg)
	pbft.State = StatePrePrepare
	pbft.PrePrepared = true

	return msg, nil
}

// ProcessPrePrepare processes a pre-prepare message (by replicas)
func (pbft *PBFT) ProcessPrePrepare(msg *PBFTMessage) error {
	pbft.mu.Lock()
	defer pbft.mu.Unlock()

	// Verify message is from primary
	if msg.NodeID != pbft.GetPrimaryNode() {
		return fmt.Errorf("pre-prepare message not from primary")
	}

	// Verify sequence and view
	if msg.Sequence != pbft.Sequence || msg.ViewID != pbft.ViewID {
		return fmt.Errorf("sequence or view mismatch")
	}

	// Verify block hash
	if msg.BlockHash != pbft.Block.Hash {
		return fmt.Errorf("block hash mismatch")
	}

	pbft.Messages = append(pbft.Messages, msg)
	pbft.State = StatePrePrepare
	pbft.PrePrepared = true

	return nil
}

// PreparePhase initiates the prepare phase
func (pbft *PBFT) PreparePhase() (*PBFTMessage, error) {
	pbft.mu.Lock()
	defer pbft.mu.Unlock()

	if !pbft.PrePrepared {
		return nil, fmt.Errorf("pre-prepare not completed")
	}

	if pbft.State != StatePrePrepare {
		return nil, fmt.Errorf("invalid state for prepare")
	}

	msg := &PBFTMessage{
		Type:      Prepare,
		BlockHash: pbft.Block.Hash,
		NodeID:    pbft.NodeID,
		Sequence:  pbft.Sequence,
		ViewID:    pbft.ViewID,
		Timestamp: time.Now(),
		Signature: pbft.signMessage(Prepare, pbft.Block.Hash),
	}

	pbft.Messages = append(pbft.Messages, msg)
	pbft.PrepareCount++
	pbft.State = StatePrepare

	return msg, nil
}

// ProcessPrepare processes a prepare message
func (pbft *PBFT) ProcessPrepare(msg *PBFTMessage) error {
	pbft.mu.Lock()
	defer pbft.mu.Unlock()

	// Verify sequence and view
	if msg.Sequence != pbft.Sequence || msg.ViewID != pbft.ViewID {
		return fmt.Errorf("sequence or view mismatch")
	}

	// Verify block hash
	if msg.BlockHash != pbft.Block.Hash {
		return fmt.Errorf("block hash mismatch")
	}

	// Check if we already have a prepare message from this node
	for _, m := range pbft.Messages {
		if m.Type == Prepare && m.NodeID == msg.NodeID {
			return nil // Already processed
		}
	}

	pbft.Messages = append(pbft.Messages, msg)
	pbft.PrepareCount++

	// Check if we have enough prepare messages (2f+1)
	if pbft.PrepareCount >= pbft.RequiredVotes {
		pbft.Prepared = true
		if pbft.State == StatePrepare {
			pbft.State = StatePrepare // Stay in prepare until commit
		}
	}

	return nil
}

// CommitPhase initiates the commit phase
func (pbft *PBFT) CommitPhase() (*PBFTMessage, error) {
	pbft.mu.Lock()
	defer pbft.mu.Unlock()

	if !pbft.Prepared {
		return nil, fmt.Errorf("prepare phase not completed")
	}

	if pbft.State != StatePrepare {
		return nil, fmt.Errorf("invalid state for commit")
	}

	msg := &PBFTMessage{
		Type:      Commit,
		BlockHash: pbft.Block.Hash,
		NodeID:    pbft.NodeID,
		Sequence:  pbft.Sequence,
		ViewID:    pbft.ViewID,
		Timestamp: time.Now(),
		Signature: pbft.signMessage(Commit, pbft.Block.Hash),
	}

	pbft.Messages = append(pbft.Messages, msg)
	pbft.CommitCount++
	pbft.State = StateCommit

	return msg, nil
}

// ProcessCommit processes a commit message
func (pbft *PBFT) ProcessCommit(msg *PBFTMessage) error {
	pbft.mu.Lock()
	defer pbft.mu.Unlock()

	// Verify sequence and view
	if msg.Sequence != pbft.Sequence || msg.ViewID != pbft.ViewID {
		return fmt.Errorf("sequence or view mismatch")
	}

	// Verify block hash
	if msg.BlockHash != pbft.Block.Hash {
		return fmt.Errorf("block hash mismatch")
	}

	// Check if we already have a commit message from this node
	for _, m := range pbft.Messages {
		if m.Type == Commit && m.NodeID == msg.NodeID {
			return nil // Already processed
		}
	}

	pbft.Messages = append(pbft.Messages, msg)
	pbft.CommitCount++

	// Check if we have enough commit messages (2f+1)
	if pbft.CommitCount >= pbft.RequiredVotes {
		pbft.Committed = true
		pbft.State = StateFinalized
	}

	return nil
}

// IsFinalized checks if the consensus is finalized
func (pbft *PBFT) IsFinalized() bool {
	pbft.mu.RLock()
	defer pbft.mu.RUnlock()
	return pbft.Committed && pbft.State == StateFinalized
}

// GetConsensusStatus returns the current consensus status
func (pbft *PBFT) GetConsensusStatus() string {
	pbft.mu.RLock()
	defer pbft.mu.RUnlock()

	return fmt.Sprintf("State: %s, PrePrepared: %v, Prepared: %v (%d/%d), Committed: %v (%d/%d)",
		pbft.State, pbft.PrePrepared, pbft.Prepared, pbft.PrepareCount, pbft.RequiredVotes,
		pbft.Committed, pbft.CommitCount, pbft.RequiredVotes)
}

// signMessage creates a simple signature for a message
func (pbft *PBFT) signMessage(msgType PBFTMessageType, blockHash string) string {
	data := fmt.Sprintf("%s:%s:%s:%d:%d", msgType, blockHash, pbft.NodeID, pbft.Sequence, pbft.ViewID)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Validate validates the PBFT consensus
func (pbft *PBFT) Validate() bool {
	pbft.mu.RLock()
	defer pbft.mu.RUnlock()

	// Check if we have enough votes
	if pbft.PrepareCount < pbft.RequiredVotes {
		return false
	}

	if pbft.CommitCount < pbft.RequiredVotes {
		return false
	}

	// Check if consensus is finalized
	return pbft.IsFinalized()
}

// CreateBlockWithPBFT creates a block using PBFT consensus
func (bc *Blockchain) CreateBlockWithPBFT(transactions []*Transaction, nodes []string, nodeID string) error {
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
		Nonce:        0, // PBFT doesn't use nonce for mining
	}

	// Calculate hash (PBFT doesn't require mining, just hash)
	newBlock.Hash = newBlock.CalculateHash()

	// Create PBFT instance
	sequence := int64(len(bc.Blocks))
	pbft := NewPBFT(nodeID, nodes, newBlock, sequence)

	// Simulate PBFT consensus process
	fmt.Printf("Starting PBFT consensus for block #%d...\n", newBlock.Index)
	fmt.Printf("  Total nodes: %d, Required votes: %d (2f+1)\n", pbft.TotalNodes, pbft.RequiredVotes)
	fmt.Printf("  Primary node: %s\n", pbft.GetPrimaryNode())

	// Phase 1: Pre-Prepare (by primary)
	if pbft.IsPrimary() {
		fmt.Println("\n  Phase 1: Pre-Prepare (Primary broadcasts block proposal)")
		msg, err := pbft.PrePreparePhase()
		if err != nil {
			return fmt.Errorf("pre-prepare phase failed: %v", err)
		}
		fmt.Printf("    Primary node sent pre-prepare message\n")
		fmt.Printf("      Block hash: %s\n", msg.BlockHash[:16]+"...")
	} else {
		// Simulate receiving pre-prepare from primary
		fmt.Println("\n  Phase 1: Pre-Prepare (Receiving from primary)")
		primaryMsg := &PBFTMessage{
			Type:      PrePrepare,
			BlockHash: newBlock.Hash,
			NodeID:    pbft.GetPrimaryNode(),
			Sequence:  sequence,
			ViewID:    0,
			Timestamp: time.Now(),
		}
		if err := pbft.ProcessPrePrepare(primaryMsg); err != nil {
			return fmt.Errorf("processing pre-prepare failed: %v", err)
		}
		fmt.Printf("    Received pre-prepare from primary\n")
	}

	// Phase 2: Prepare (all nodes)
	fmt.Println("\n  Phase 2: Prepare (Nodes validate and broadcast prepare)")
	if _, err := pbft.PreparePhase(); err != nil {
		return fmt.Errorf("prepare phase failed: %v", err)
	}
	fmt.Printf("    Node %s sent prepare message\n", nodeID[:16]+"...")

	// Simulate receiving prepare messages from other nodes
	for i, node := range nodes {
		if node != nodeID {
			msg := &PBFTMessage{
				Type:      Prepare,
				BlockHash: newBlock.Hash,
				NodeID:    node,
				Sequence:  sequence,
				ViewID:    0,
				Timestamp: time.Now(),
			}
			pbft.ProcessPrepare(msg)
			if i < 3 { // Show first 3 for clarity
				fmt.Printf("    Received prepare from node %s\n", node[:16]+"...")
			}
		}
	}
	fmt.Printf("    Total prepare messages: %d/%d\n", pbft.PrepareCount, pbft.RequiredVotes)

	if !pbft.Prepared {
		return fmt.Errorf("failed to reach prepare quorum")
	}
	fmt.Println("    Prepare phase completed (quorum reached)")

	// Phase 3: Commit (all nodes)
	fmt.Println("\n  Phase 3: Commit (Nodes broadcast commit)")
	if _, err := pbft.CommitPhase(); err != nil {
		return fmt.Errorf("commit phase failed: %v", err)
	}
	fmt.Printf("    Node %s sent commit message\n", nodeID[:16]+"...")

	// Simulate receiving commit messages from other nodes
	for i, node := range nodes {
		if node != nodeID {
			msg := &PBFTMessage{
				Type:      Commit,
				BlockHash: newBlock.Hash,
				NodeID:    node,
				Sequence:  sequence,
				ViewID:    0,
				Timestamp: time.Now(),
			}
			pbft.ProcessCommit(msg)
			if i < 3 { // Show first 3 for clarity
				fmt.Printf("    Received commit from node %s\n", node[:16]+"...")
			}
		}
	}
	fmt.Printf("    Total commit messages: %d/%d\n", pbft.CommitCount, pbft.RequiredVotes)

	if !pbft.IsFinalized() {
		return fmt.Errorf("failed to reach commit quorum")
	}
	fmt.Println("    Commit phase completed (quorum reached)")
	fmt.Println("    Block finalized with PBFT consensus!")

	// Validate consensus
	if !pbft.Validate() {
		return fmt.Errorf("PBFT consensus validation failed")
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	fmt.Printf("\nBlock #%d added to the blockchain using PBFT!\n", newBlock.Index)
	fmt.Printf("  Byzantine fault tolerance: Can tolerate %d faulty nodes\n\n", (pbft.TotalNodes-1)/3)

	return nil
}
