package main

import "fmt"

// GetBalance calculates the balance of an address by scanning all transactions
func (bc *Blockchain) GetBalance(address string) float64 {
	balance := 0.0

	// Scan all blocks
	for _, block := range bc.Blocks {
		// Scan all transactions in the block
		for _, tx := range block.Transactions {
			// Skip genesis transaction
			if tx.From == "" && tx.To == "Genesis" {
				continue
			}

			// Subtract if address is sender
			if tx.From == address {
				balance -= tx.Amount
			}

			// Add if address is receiver
			if tx.To == address {
				balance += tx.Amount
			}
		}
	}

	return balance
}

// ValidateTransaction checks if a transaction is valid (sufficient balance)
func (bc *Blockchain) ValidateTransaction(tx *Transaction) error {
	// Skip validation for genesis-like transactions
	if tx.From == "" {
		return nil
	}

	balance := bc.GetBalance(tx.From)
	if balance < tx.Amount {
		return fmt.Errorf("insufficient balance: address %s has %.2f, trying to send %.2f", tx.From, balance, tx.Amount)
	}

	return nil
}

// AddCoinbaseTransaction creates a coinbase transaction to give initial balance
func (bc *Blockchain) AddCoinbaseTransaction(to string, amount float64) *Transaction {
	// Coinbase transaction has empty From address
	return NewTransaction("", to, amount)
}
