package main

import (
	"encoding/hex"
	"fmt"
	"sync"
)

// Mempool represents a transaction pool for pending transactions
type Mempool struct {
	transactions map[string]*Transaction // Map by transaction hash
	mu           sync.RWMutex
}

// NewMempool creates a new mempool
func NewMempool() *Mempool {
	return &Mempool{
		transactions: make(map[string]*Transaction),
	}
}

// AddTransaction adds a transaction to the mempool
func (mp *Mempool) AddTransaction(tx *Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	txHash := hex.EncodeToString(tx.Hash())

	// Check if transaction already exists
	if _, exists := mp.transactions[txHash]; exists {
		return fmt.Errorf("transaction already exists in mempool")
	}

	mp.transactions[txHash] = tx
	return nil
}

// GetTransaction retrieves a transaction by hash
func (mp *Mempool) GetTransaction(txHash string) (*Transaction, bool) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	tx, exists := mp.transactions[txHash]
	return tx, exists
}

// GetAllTransactions returns all transactions in the mempool
func (mp *Mempool) GetAllTransactions() []*Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	transactions := make([]*Transaction, 0, len(mp.transactions))
	for _, tx := range mp.transactions {
		transactions = append(transactions, tx)
	}
	return transactions
}

// RemoveTransaction removes a transaction from the mempool
func (mp *Mempool) RemoveTransaction(txHash string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	delete(mp.transactions, txHash)
}

// RemoveTransactions removes multiple transactions from the mempool
func (mp *Mempool) RemoveTransactions(txHashes []string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	for _, txHash := range txHashes {
		delete(mp.transactions, txHash)
	}
}

// Size returns the number of transactions in the mempool
func (mp *Mempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return len(mp.transactions)
}

// Clear removes all transactions from the mempool
func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.transactions = make(map[string]*Transaction)
}

// GetTransactionsForBlock returns up to maxTransactions transactions for a new block
func (mp *Mempool) GetTransactionsForBlock(maxTransactions int) []*Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	transactions := make([]*Transaction, 0, maxTransactions)
	count := 0
	for _, tx := range mp.transactions {
		if count >= maxTransactions {
			break
		}
		transactions = append(transactions, tx)
		count++
	}
	return transactions
}
