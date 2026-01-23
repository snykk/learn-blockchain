package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// BridgeDirection represents the direction of the bridge transfer
type BridgeDirection string

const (
	BridgeDirectionAToB BridgeDirection = "a_to_b" // Chain A to Chain B
	BridgeDirectionBToA BridgeDirection = "b_to_a" // Chain B to Chain A
)

// BridgeStatus represents the status of a bridge transaction
type BridgeStatus string

const (
	BridgeStatusPending   BridgeStatus = "pending"   // Waiting for validator approvals
	BridgeStatusApproved  BridgeStatus = "approved"  // Approved by validators
	BridgeStatusCompleted BridgeStatus = "completed"  // Transfer completed
	BridgeStatusRejected  BridgeStatus = "rejected"  // Rejected by validators
)

// BridgeTransaction represents a cross-chain transfer
type BridgeTransaction struct {
	TxID           string          `json:"tx_id"`
	FromChain      string          `json:"from_chain"`
	ToChain        string          `json:"to_chain"`
	FromAddress    string          `json:"from_address"`
	ToAddress      string          `json:"to_address"`
	Amount         float64         `json:"amount"`
	Token          string          `json:"token"`
	Status         BridgeStatus    `json:"status"`
	Direction      BridgeDirection `json:"direction"`
	Timestamp      time.Time       `json:"timestamp"`
	Approvals      int             `json:"approvals"`
	RequiredSigs   int             `json:"required_sigs"`
	Signatures     []string        `json:"signatures"`
	LockTxHash     string          `json:"lock_tx_hash"`     // Tx hash on source chain
	UnlockTxHash   string          `json:"unlock_tx_hash"`   // Tx hash on destination chain
}

// BridgeEvent represents an event emitted by the bridge
type BridgeEvent struct {
	EventType   string    `json:"event_type"`   // lock, unlock, approval
	Chain       string    `json:"chain"`
	TxHash      string    `json:"tx_hash"`
	Timestamp   time.Time `json:"timestamp"`
	Data        string    `json:"data"`
}

// Validator represents a bridge validator
type Validator struct {
	ID          string  `json:"id"`
	Address     string  `json:"address"`
	Stake       float64 `json:"stake"`
	IsActive    bool    `json:"is_active"`
	VotingPower int     `json:"voting_power"`
}

// Bridge represents a cross-chain bridge between two blockchains
type Bridge struct {
	BridgeID       string
	ChainA         *Blockchain
	ChainB         *Blockchain
	ChainAName     string
	ChainBName     string
	Validators     []*Validator
	RequiredSigs   int
	PendingTxs     map[string]*BridgeTransaction
	CompletedTxs   map[string]*BridgeTransaction
	Events         []*BridgeEvent
	mu             sync.RWMutex
	MinAmount      float64
	MaxAmount      float64
	Fee            float64
	RelayerAddress string
}

// BridgeManager manages multiple bridges
type BridgeManager struct {
	Bridges    map[string]*Bridge
	mu         sync.RWMutex
	Blockchain *Blockchain
}

// NewBridgeManager creates a new bridge manager
func NewBridgeManager(bc *Blockchain) *BridgeManager {
	return &BridgeManager{
		Bridges:    make(map[string]*Bridge),
		Blockchain: bc,
	}
}

// NewBridge creates a new cross-chain bridge
func NewBridge(bridgeID string, chainA, chainB *Blockchain, chainAName, chainBName string, requiredSigs int) *Bridge {
	bridge := &Bridge{
		BridgeID:     bridgeID,
		ChainA:       chainA,
		ChainB:       chainB,
		ChainAName:   chainAName,
		ChainBName:   chainBName,
		Validators:   make([]*Validator, 0),
		RequiredSigs: requiredSigs,
		PendingTxs:   make(map[string]*BridgeTransaction),
		CompletedTxs: make(map[string]*BridgeTransaction),
		Events:       make([]*BridgeEvent, 0),
		MinAmount:    0.1,
		MaxAmount:    10000.0,
		Fee:          0.01, // 1% bridge fee
		RelayerAddress: "relayer_" + bridgeID,
	}

	return bridge
}

// AddValidator adds a validator to the bridge
func (b *Bridge) AddValidator(id, address string, stake float64, votingPower int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	validator := &Validator{
		ID:          id,
		Address:     address,
		Stake:       stake,
		IsActive:    true,
		VotingPower: votingPower,
	}

	b.Validators = append(b.Validators, validator)

	fmt.Printf("\n[Bridge Validator Added]\n")
	fmt.Printf("  Bridge: %s\n", b.BridgeID[:16]+"...")
	fmt.Printf("  Validator ID: %s\n", validator.ID[:16]+"...")
	fmt.Printf("  Address: %s\n", validator.Address[:16]+"...")
	fmt.Printf("  Stake: %.2f\n", validator.Stake)
	fmt.Printf("  Voting Power: %d\n", validator.VotingPower)
}

// LockFunds locks funds on the source chain (Chain A -> Chain B)
func (b *Bridge) LockFunds(fromAddress, toAddress string, amount float64, token string) (*BridgeTransaction, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Validate amount
	if amount < b.MinAmount {
		return nil, fmt.Errorf("amount below minimum: %.4f < %.4f", amount, b.MinAmount)
	}
	if amount > b.MaxAmount {
		return nil, fmt.Errorf("amount above maximum: %.4f > %.4f", amount, b.MaxAmount)
	}

	// Calculate fee
	fee := amount * b.Fee
	totalAmount := amount + fee

	// Check balance on Chain A
	balance := b.ChainA.GetBalance(fromAddress)
	if balance < totalAmount {
		return nil, fmt.Errorf("insufficient balance on %s: %.2f < %.2f", b.ChainAName, balance, totalAmount)
	}

	// Lock funds on Chain A (create lock transaction)
	// In a real implementation, this would call a bridge smart contract
	lockTxHash := generateLockTxHash(fromAddress, toAddress, amount, time.Now())

	// Create bridge transaction
	txID := generateBridgeTxID(lockTxHash, b.ChainAName, b.ChainBName)

	bridgeTx := &BridgeTransaction{
		TxID:         txID,
		FromChain:    b.ChainAName,
		ToChain:      b.ChainBName,
		FromAddress:  fromAddress,
		ToAddress:    toAddress,
		Amount:       amount,
		Token:        token,
		Status:       BridgeStatusPending,
		Direction:    BridgeDirectionAToB,
		Timestamp:    time.Now(),
		Approvals:    0,
		RequiredSigs: b.RequiredSigs,
		Signatures:   make([]string, 0),
		LockTxHash:   lockTxHash,
	}

	b.PendingTxs[txID] = bridgeTx

	// Emit lock event
	b.emitEvent("lock", b.ChainAName, lockTxHash, fmt.Sprintf("%s->%s: %.4f %s", b.ChainAName, b.ChainBName, amount, token))

	fmt.Printf("\n=== Cross-Chain Bridge: Lock Funds ===\n")
	fmt.Printf("Bridge: %s\n", b.BridgeID[:16]+"...")
	fmt.Printf("Direction: %s → %s\n", b.ChainAName, b.ChainBName)
	fmt.Printf("From: %s\n", fromAddress[:16]+"...")
	fmt.Printf("To: %s\n", toAddress[:16]+"...")
	fmt.Printf("Amount: %.4f %s\n", amount, token)
	fmt.Printf("Fee: %.4f %s (%.1f%%)\n", fee, token, b.Fee*100)
	fmt.Printf("Lock Tx Hash: %s\n", lockTxHash[:16]+"...")
	fmt.Printf("Status: %s\n", bridgeTx.Status)
	fmt.Printf("Required Signatures: %d/%d\n", 0, b.RequiredSigs)

	return bridgeTx, nil
}

// UnlockFunds unlocks/mints funds on the destination chain
func (b *Bridge) UnlockFunds(bridgeTx *BridgeTransaction) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Verify transaction is approved
	if bridgeTx.Status != BridgeStatusApproved {
		return fmt.Errorf("transaction not approved: %s", bridgeTx.Status)
	}

	// Unlock/mint funds on Chain B
	// In a real implementation, this would call a bridge smart contract on Chain B
	unlockTxHash := generateUnlockTxHash(bridgeTx.ToAddress, bridgeTx.Amount, time.Now())
	bridgeTx.UnlockTxHash = unlockTxHash

	// Add coinbase transaction to Chain B
	coinbaseTx := NewTransaction(b.RelayerAddress, bridgeTx.ToAddress, bridgeTx.Amount)
	b.ChainB.AddBlock([]*Transaction{coinbaseTx})

	// Move to completed
	delete(b.PendingTxs, bridgeTx.TxID)
	bridgeTx.Status = BridgeStatusCompleted
	b.CompletedTxs[bridgeTx.TxID] = bridgeTx

	// Emit unlock event
	b.emitEvent("unlock", b.ChainBName, unlockTxHash, fmt.Sprintf("Minted %.4f %s to %s", bridgeTx.Amount, bridgeTx.Token, bridgeTx.ToAddress[:16]+"..."))

	fmt.Printf("\n=== Cross-Chain Bridge: Unlock Funds ===\n")
	fmt.Printf("Bridge: %s\n", b.BridgeID[:16]+"...")
	fmt.Printf("Direction: %s → %s\n", b.ChainAName, b.ChainBName)
	fmt.Printf("To: %s\n", bridgeTx.ToAddress[:16]+"...")
	fmt.Printf("Amount: %.4f %s\n", bridgeTx.Amount, bridgeTx.Token)
	fmt.Printf("Unlock Tx Hash: %s\n", unlockTxHash[:16]+"...")
	fmt.Printf("✓ Funds successfully transferred to %s\n", b.ChainBName)

	return nil
}

// ApproveTransaction approves a bridge transaction (by validator)
func (b *Bridge) ApproveTransaction(txID, validatorID string, signature string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	bridgeTx, exists := b.PendingTxs[txID]
	if !exists {
		return fmt.Errorf("transaction not found: %s", txID)
	}

	// Check if already approved by this validator
	for _, sig := range bridgeTx.Signatures {
		if sig == signature {
			return fmt.Errorf("already approved by validator")
		}
	}

	// Add signature
	bridgeTx.Signatures = append(bridgeTx.Signatures, signature)
	bridgeTx.Approvals++

	fmt.Printf("\n[Bridge Transaction Approved]\n")
	fmt.Printf("  Tx ID: %s\n", txID[:16]+"...")
	fmt.Printf("  Validator: %s\n", validatorID[:16]+"...")
	fmt.Printf("  Approvals: %d/%d\n", bridgeTx.Approvals, bridgeTx.RequiredSigs)

	// Check if we have enough approvals
	if bridgeTx.Approvals >= bridgeTx.RequiredSigs {
		bridgeTx.Status = BridgeStatusApproved
		fmt.Printf("  ✓ Transaction approved by validators!\n")
	}

	// Emit approval event
	b.emitEvent("approval", b.ChainAName, txID, fmt.Sprintf("Validator %s approved", validatorID[:16]+"..."))

	return nil
}

// ReverseTransfer reverses direction (Chain B -> Chain A)
func (b *Bridge) ReverseTransfer(fromAddress, toAddress string, amount float64, token string) (*BridgeTransaction, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Lock funds on Chain B
	balance := b.ChainB.GetBalance(fromAddress)
	totalAmount := amount + (amount * b.Fee)

	if balance < totalAmount {
		return nil, fmt.Errorf("insufficient balance on %s: %.2f < %.2f", b.ChainBName, balance, totalAmount)
	}

	lockTxHash := generateLockTxHash(fromAddress, toAddress, amount, time.Now())
	txID := generateBridgeTxID(lockTxHash, b.ChainBName, b.ChainAName)

	bridgeTx := &BridgeTransaction{
		TxID:         txID,
		FromChain:    b.ChainBName,
		ToChain:      b.ChainAName,
		FromAddress:  fromAddress,
		ToAddress:    toAddress,
		Amount:       amount,
		Token:        token,
		Status:       BridgeStatusPending,
		Direction:    BridgeDirectionBToA,
		Timestamp:    time.Now(),
		Approvals:    0,
		RequiredSigs: b.RequiredSigs,
		Signatures:   make([]string, 0),
		LockTxHash:   lockTxHash,
	}

	b.PendingTxs[txID] = bridgeTx

	b.emitEvent("lock", b.ChainBName, lockTxHash, fmt.Sprintf("%s->%s: %.4f %s", b.ChainBName, b.ChainAName, amount, token))

	fmt.Printf("\n=== Cross-Chain Bridge: Reverse Transfer ===\n")
	fmt.Printf("Bridge: %s\n", b.BridgeID[:16]+"...")
	fmt.Printf("Direction: %s → %s (REVERSE)\n", b.ChainBName, b.ChainAName)
	fmt.Printf("From: %s\n", fromAddress[:16]+"...")
	fmt.Printf("To: %s\n", toAddress[:16]+"...")
	fmt.Printf("Amount: %.4f %s\n", amount, token)

	return bridgeTx, nil
}

// GetBridgeStatistics returns bridge statistics
func (b *Bridge) GetBridgeStatistics() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	totalVolume := 0.0
	pendingVolume := 0.0

	for _, tx := range b.CompletedTxs {
		totalVolume += tx.Amount
	}

	for _, tx := range b.PendingTxs {
		pendingVolume += tx.Amount
	}

	return map[string]interface{}{
		"bridge_id":         b.BridgeID[:16] + "...",
		"chain_a":           b.ChainAName,
		"chain_b":           b.ChainBName,
		"validators":        len(b.Validators),
		"required_sigs":     b.RequiredSigs,
		"pending_txs":       len(b.PendingTxs),
		"completed_txs":     len(b.CompletedTxs),
		"total_volume":      totalVolume,
		"pending_volume":    pendingVolume,
		"fee":               b.Fee,
		"events":            len(b.Events),
	}
}

// GetTransaction retrieves a transaction by ID
func (b *Bridge) GetTransaction(txID string) (*BridgeTransaction, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Check pending
	if tx, exists := b.PendingTxs[txID]; exists {
		return tx, nil
	}

	// Check completed
	if tx, exists := b.CompletedTxs[txID]; exists {
		return tx, nil
	}

	return nil, fmt.Errorf("transaction not found: %s", txID)
}

// emitEvent emits a bridge event
func (b *Bridge) emitEvent(eventType, chain, txHash, data string) {
	event := &BridgeEvent{
		EventType: eventType,
		Chain:     chain,
		TxHash:    txHash,
		Timestamp: time.Now(),
		Data:      data,
	}

	b.Events = append(b.Events, event)
}

// Helper functions

func generateBridgeTxID(lockTxHash, fromChain, toChain string) string {
	data := fmt.Sprintf("bridge:%s:%s:%s:%d", lockTxHash, fromChain, toChain, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateLockTxHash(from, to string, amount float64, timestamp time.Time) string {
	data := fmt.Sprintf("lock:%s:%s:%.4f:%d", from, to, amount, timestamp.UnixNano())
	hash := sha256.Sum256([]byte(data))
	return "L" + hex.EncodeToString(hash[:])[:40]
}

func generateUnlockTxHash(to string, amount float64, timestamp time.Time) string {
	data := fmt.Sprintf("unlock:%s:%.4f:%d", to, amount, timestamp.UnixNano())
	hash := sha256.Sum256([]byte(data))
	return "U" + hex.EncodeToString(hash[:])[:40]
}
