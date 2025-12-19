package main

import "fmt"

const (
	// BlockReward is the reward given to miners/validators for creating a block
	BlockReward = 50.0
	// InitialBlockReward is the reward for the genesis block
	InitialBlockReward = 100.0
)

// BlockRewardTransaction creates a block reward transaction for the miner/validator
func NewBlockRewardTransaction(minerAddress string, isGenesis bool) *Transaction {
	reward := BlockReward
	if isGenesis {
		reward = InitialBlockReward
	}

	// Block reward transaction has empty From address (new coins created)
	return NewTransaction("", minerAddress, reward)
}

// GetMinerRewards calculates total rewards earned by a miner/validator
func (bc *Blockchain) GetMinerRewards(minerAddress string) float64 {
	rewards := 0.0

	for i, block := range bc.Blocks {
		// Count block reward transactions
		for _, tx := range block.Transactions {
			if tx.From == "" && tx.To == minerAddress {
				rewards += tx.Amount
			}
		}

		// Count transaction fees (simplified: first fee per block)
		if i > 0 {
			for _, tx := range block.Transactions {
				if tx.Fee > 0 {
					rewards += tx.Fee
					break
				}
			}
		}
	}

	return rewards
}

// GetTotalBalance returns the total balance including rewards
func (bc *Blockchain) GetTotalBalance(address string) float64 {
	balance := bc.GetBalance(address)
	rewards := bc.GetMinerRewards(address)
	return balance + rewards
}

// CalculateTotalFees calculates total fees from transactions in a block
func CalculateTotalFees(transactions []*Transaction) float64 {
	totalFees := 0.0
	for _, tx := range transactions {
		totalFees += tx.Fee
	}
	return totalFees
}

// FormatRewardInfo returns a formatted string for block reward info
func FormatRewardInfo(minerAddress string, blockReward, totalFees float64) string {
	if totalFees > 0 {
		return fmt.Sprintf("Miner: %s, Block Reward: %.2f, Fees: %.2f, Total: %.2f",
			minerAddress[:16]+"...", blockReward, totalFees, blockReward+totalFees)
	}
	return fmt.Sprintf("Miner: %s, Block Reward: %.2f", minerAddress[:16]+"...", blockReward)
}
