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

			// Subtract if address is sender (amount + fee)
			if tx.From == address {
				balance -= tx.Amount
				balance -= tx.Fee // Subtract transaction fee
			}

			// Add if address is receiver
			if tx.To == address {
				balance += tx.Amount
			}

			// Add if address is miner (from block rewards)
			// Block rewards are handled separately in GetMinerRewards
		}
	}

	return balance
}

// ValidateTransaction checks if a transaction is valid (sufficient balance including fee)
func (bc *Blockchain) ValidateTransaction(tx *Transaction) error {
	// Skip validation for genesis-like transactions
	if tx.From == "" {
		return nil
	}

	balance := bc.GetBalance(tx.From)
	totalCost := tx.TotalCost() // Amount + Fee
	if balance < totalCost {
		return fmt.Errorf("insufficient balance: address %s has %.2f, trying to send %.2f (amount) + %.2f (fee) = %.2f (total)",
			tx.From, balance, tx.Amount, tx.Fee, totalCost)
	}

	return nil
}

// AddCoinbaseTransaction creates a coinbase transaction to give initial balance
func (bc *Blockchain) AddCoinbaseTransaction(to string, amount float64) *Transaction {
	// Coinbase transaction has empty From address
	return NewTransaction("", to, amount)
}
