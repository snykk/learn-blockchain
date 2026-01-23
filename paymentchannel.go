package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// ChannelState represents the state of a payment channel
type ChannelState struct {
	ChannelID      string    `json:"channel_id"`
	Participant1   string    `json:"participant1"`
	Participant2   string    `json:"participant2"`
	Balance1       float64   `json:"balance1"`
	Balance2       float64   `json:"balance2"`
	SequenceNumber int64     `json:"sequence_number"`
	Nonce          int64     `json:"nonce"`
	Timestamp      time.Time `json:"timestamp"`
	IsClosed       bool      `json:"is_closed"`
	ClosingTxHash  string    `json:"closing_tx_hash,omitempty"`
}

// ChannelSignature represents a signed channel state
type ChannelSignature struct {
	State      *ChannelState `json:"state"`
	Signature1 string        `json:"signature1"` // From participant1
	Signature2 string        `json:"signature2"` // From participant2
}

// PaymentChannel represents a payment channel (Layer 2 solution)
type PaymentChannel struct {
	State           *ChannelState
	InitialState    *ChannelState
	DepositAmount   float64
	MultiSigAddress string
	Timeout         time.Duration
	CreatedAt       time.Time
	LastUpdate      time.Time
	mu              sync.RWMutex
	Blockchain      *Blockchain
	PendingUpdates  []*ChannelState
	UpdateHistory   []*ChannelState
}

// ChannelManager manages multiple payment channels
type ChannelManager struct {
	Channels   map[string]*PaymentChannel
	mu         sync.RWMutex
	Blockchain *Blockchain
}

// NewChannelManager creates a new channel manager
func NewChannelManager(bc *Blockchain) *ChannelManager {
	return &ChannelManager{
		Channels:   make(map[string]*PaymentChannel),
		Blockchain: bc,
	}
}

// CreateChannel creates a new payment channel between two parties
func (cm *ChannelManager) CreateChannel(participant1, participant2 string, deposit1, deposit2 float64, timeout time.Duration) (*PaymentChannel, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Validate deposits
	if deposit1 <= 0 || deposit2 <= 0 {
		return nil, fmt.Errorf("deposits must be positive")
	}

	// Validate balances
	balance1 := cm.Blockchain.GetBalance(participant1)
	balance2 := cm.Blockchain.GetBalance(participant2)

	if balance1 < deposit1 {
		return nil, fmt.Errorf("participant1 has insufficient balance: %.2f < %.2f", balance1, deposit1)
	}
	if balance2 < deposit2 {
		return nil, fmt.Errorf("participant2 has insufficient balance: %.2f < %.2f", balance2, deposit2)
	}

	// Generate channel ID
	channelID := generateChannelID(participant1, participant2, time.Now())

	// Create initial state
	initialState := &ChannelState{
		ChannelID:      channelID,
		Participant1:   participant1,
		Participant2:   participant2,
		Balance1:       deposit1,
		Balance2:       deposit2,
		SequenceNumber: 0,
		Nonce:          0,
		Timestamp:      time.Now(),
		IsClosed:       false,
	}

	// Create multisig address (simplified - in reality this would be a proper multisig)
	multiSigAddress := generateMultisigAddress(participant1, participant2, channelID)

	channel := &PaymentChannel{
		State:           initialState,
		InitialState:    initialState,
		DepositAmount:   deposit1 + deposit2,
		MultiSigAddress: multiSigAddress,
		Timeout:         timeout,
		CreatedAt:       time.Now(),
		LastUpdate:      time.Now(),
		Blockchain:      cm.Blockchain,
		PendingUpdates:  make([]*ChannelState, 0),
		UpdateHistory:   []*ChannelState{initialState},
	}

	cm.Channels[channelID] = channel

	fmt.Printf("\n=== Payment Channel Created ===\n")
	fmt.Printf("Channel ID: %s\n", channelID[:16]+"...")
	fmt.Printf("Participants: %s ↔ %s\n", participant1[:16]+"...", participant2[:16]+"...")
	fmt.Printf("Initial Balances: %.2f / %.2f\n", deposit1, deposit2)
	fmt.Printf("Total Deposit: %.2f\n", deposit1+deposit2)
	fmt.Printf("Timeout: %v\n", timeout)
	fmt.Printf("Multisig Address: %s\n", multiSigAddress[:16]+"...")

	return channel, nil
}

// GetChannel retrieves a channel by ID
func (cm *ChannelManager) GetChannel(channelID string) (*PaymentChannel, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	channel, exists := cm.Channels[channelID]
	if !exists {
		return nil, fmt.Errorf("channel not found: %s", channelID)
	}

	return channel, nil
}

// UpdateState proposes a new state for the channel
func (pc *PaymentChannel) UpdateState(newBalance1, newBalance2 float64) (*ChannelState, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.State.IsClosed {
		return nil, fmt.Errorf("channel is closed")
	}

	// Validate balances
	total := pc.State.Balance1 + pc.State.Balance2
	if newBalance1+newBalance2 != total {
		return nil, fmt.Errorf("total balance must remain constant: %.2f != %.2f", newBalance1+newBalance2, total)
	}

	if newBalance1 < 0 || newBalance2 < 0 {
		return nil, fmt.Errorf("balances cannot be negative")
	}

	// Create new state
	newState := &ChannelState{
		ChannelID:      pc.State.ChannelID,
		Participant1:   pc.State.Participant1,
		Participant2:   pc.State.Participant2,
		Balance1:       newBalance1,
		Balance2:       newBalance2,
		SequenceNumber: pc.State.SequenceNumber + 1,
		Nonce:          pc.State.Nonce + 1,
		Timestamp:      time.Now(),
		IsClosed:       false,
	}

	// Add to pending updates (waiting for signatures)
	pc.PendingUpdates = append(pc.PendingUpdates, newState)

	fmt.Printf("\n[Channel Update Proposed]\n")
	fmt.Printf("  New State: %.2f ↔ %.2f\n", newBalance1, newBalance2)
	fmt.Printf("  Sequence Number: %d\n", newState.SequenceNumber)

	return newState, nil
}

// SignState signs a channel state (in reality, this would use wallet signing)
func (pc *PaymentChannel) SignState(state *ChannelState, signer string) (string, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	// Verify signer is a participant
	if signer != state.Participant1 && signer != state.Participant2 {
		return "", fmt.Errorf("signer is not a participant")
	}

	// Create signature (simplified - in reality this would use ECDSA)
	signature := signChannelState(state, signer)

	return signature, nil
}

// CommitState commits a signed state to the channel
func (pc *PaymentChannel) CommitState(signedState *ChannelSignature) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.State.IsClosed {
		return fmt.Errorf("channel is closed")
	}

	// Verify signatures (simplified)
	if signedState.Signature1 == "" || signedState.Signature2 == "" {
		return fmt.Errorf("both signatures are required")
	}

	// Verify sequence number
	if signedState.State.SequenceNumber <= pc.State.SequenceNumber {
		return fmt.Errorf("invalid sequence number")
	}

	// Update state
	pc.State = signedState.State
	pc.LastUpdate = time.Now()

	// Add to history
	pc.UpdateHistory = append(pc.UpdateHistory, signedState.State)

	// Clear pending updates
	pc.PendingUpdates = make([]*ChannelState, 0)

	fmt.Printf("\n[Channel State Committed]\n")
	fmt.Printf("  Sequence: %d\n", pc.State.SequenceNumber)
	fmt.Printf("  Balances: %.2f ↔ %.2f\n", pc.State.Balance1, pc.State.Balance2)

	return nil
}

// CloseChannel closes the payment channel and settles on blockchain
func (pc *PaymentChannel) CloseChannel(finalState *ChannelSignature) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.State.IsClosed {
		return fmt.Errorf("channel is already closed")
	}

	// Verify final state
	if finalState.State.SequenceNumber != pc.State.SequenceNumber {
		return fmt.Errorf("final state sequence number mismatch")
	}

	// Mark as closed
	finalState.State.IsClosed = true
	pc.State = finalState.State
	pc.State.ClosingTxHash = generateClosingTxHash(finalState.State)

	fmt.Printf("\n=== Payment Channel Closing ===\n")
	fmt.Printf("Channel ID: %s\n", pc.State.ChannelID[:16]+"...")
	fmt.Printf("Final Balances: %.2f ↔ %.2f\n", pc.State.Balance1, pc.State.Balance2)
	fmt.Printf("Closing Transaction Hash: %s\n", pc.State.ClosingTxHash[:16]+"...")

	// In a real implementation, this would create a closing transaction on the blockchain
	// For now, we just simulate it
	fmt.Printf("✓ Channel closed successfully\n")
	fmt.Printf("✓ Funds settled on blockchain\n")

	return nil
}

// GetStatus returns the current status of the channel
func (pc *PaymentChannel) GetStatus() string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	duration := time.Since(pc.CreatedAt)
	transactions := len(pc.UpdateHistory) - 1 // Exclude initial state

	return fmt.Sprintf(
		"Channel: %s...\n  Status: %s\n  Balances: %.2f / %.2f\n  Transactions: %d\n  Duration: %v\n  Sequence: %d",
		pc.State.ChannelID[:16],
		map[bool]string{true: "Closed", false: "Open"}[pc.State.IsClosed],
		pc.State.Balance1,
		pc.State.Balance2,
		transactions,
		duration.Round(time.Second),
		pc.State.SequenceNumber,
	)
}

// GetChannelStatistics returns statistics about all channels
func (cm *ChannelManager) GetChannelStatistics() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	totalChannels := len(cm.Channels)
	openChannels := 0
	closedChannels := 0
	totalTransactions := 0
	totalVolume := 0.0

	for _, channel := range cm.Channels {
		if channel.State.IsClosed {
			closedChannels++
		} else {
			openChannels++
		}
		totalTransactions += len(channel.UpdateHistory) - 1
		totalVolume += channel.DepositAmount
	}

	return map[string]interface{}{
		"total_channels":     totalChannels,
		"open_channels":      openChannels,
		"closed_channels":    closedChannels,
		"total_transactions": totalTransactions,
		"total_volume":       totalVolume,
	}
}

// MicroPayment performs a micropayment through the channel
func (pc *PaymentChannel) MicroPayment(sender string, amount float64) (*ChannelState, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.State.IsClosed {
		return nil, fmt.Errorf("channel is closed")
	}

	var newBalance1, newBalance2 float64

	switch sender {
	case pc.State.Participant1:
		// Participant1 sends to Participant2
		if pc.State.Balance1 < amount {
			return nil, fmt.Errorf("insufficient balance: %.2f < %.2f", pc.State.Balance1, amount)
		}
		newBalance1 = pc.State.Balance1 - amount
		newBalance2 = pc.State.Balance2 + amount
	case pc.State.Participant2:
		// Participant2 sends to Participant1
		if pc.State.Balance2 < amount {
			return nil, fmt.Errorf("insufficient balance: %.2f < %.2f", pc.State.Balance2, amount)
		}
		newBalance1 = pc.State.Balance1 + amount
		newBalance2 = pc.State.Balance2 - amount
	default:
		return nil, fmt.Errorf("sender is not a participant")
	}

	// Create new state
	newState := &ChannelState{
		ChannelID:      pc.State.ChannelID,
		Participant1:   pc.State.Participant1,
		Participant2:   pc.State.Participant2,
		Balance1:       newBalance1,
		Balance2:       newBalance2,
		SequenceNumber: pc.State.SequenceNumber + 1,
		Nonce:          pc.State.Nonce + 1,
		Timestamp:      time.Now(),
		IsClosed:       false,
	}

	fmt.Printf("\n[Micropayment via Channel]\n")
	fmt.Printf("  From: %s\n", sender[:16]+"...")
	fmt.Printf("  Amount: %.4f\n", amount)
	fmt.Printf("  New Balances: %.2f ↔ %.2f\n", newBalance1, newBalance2)

	return newState, nil
}

// Helper functions

func generateChannelID(participant1, participant2 string, timestamp time.Time) string {
	data := fmt.Sprintf("%s:%s:%d", participant1, participant2, timestamp.UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateMultisigAddress(participant1, participant2, channelID string) string {
	data := fmt.Sprintf("multisig:%s:%s:%s", participant1, participant2, channelID)
	hash := sha256.Sum256([]byte(data))
	return "M" + hex.EncodeToString(hash[:])[:40]
}

func signChannelState(state *ChannelState, signer string) string {
	data := fmt.Sprintf("%s:%.2f:%.2f:%d:%s", state.ChannelID, state.Balance1, state.Balance2, state.SequenceNumber, signer)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateClosingTxHash(state *ChannelState) string {
	data := fmt.Sprintf("closing:%s:%.2f:%.2f:%d", state.ChannelID, state.Balance1, state.Balance2, state.SequenceNumber)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
