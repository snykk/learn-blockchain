package main

import (
	"encoding/hex"
	"fmt"
	"time"
)

// Blockchain represents a blockchain
type Blockchain struct {
	Blocks  []*Block
	Mempool *Mempool
}

// NewBlockchain creates a new blockchain with genesis block
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:  []*Block{},
		Mempool: NewMempool(),
	}
	bc.CreateGenesisBlock()
	return bc
}

// CreateGenesisBlock creates the first block in the blockchain
func (bc *Blockchain) CreateGenesisBlock() {
	// Create genesis transaction
	genesisTx := NewTransaction("", "Genesis", 0)
	transactions := []*Transaction{genesisTx}

	// Create Merkle tree
	merkleTree := NewMerkleTree(transactions)
	merkleRoot := merkleTree.GetRootHash()

	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: transactions,
		MerkleRoot:   merkleRoot,
		PreviousHash: "0",
		Nonce:        0,
	}

	// Mine the genesis block
	pow := NewProofOfWork(genesisBlock)
	nonce, hash := pow.Run()
	genesisBlock.Nonce = nonce
	genesisBlock.Hash = hash

	bc.Blocks = append(bc.Blocks, genesisBlock)
	fmt.Println("Genesis block created and mined!")
}

// AddTransactionToMempool adds a transaction to the mempool
func (bc *Blockchain) AddTransactionToMempool(tx *Transaction) error {
	// Validate transaction first
	if err := bc.ValidateTransaction(tx); err != nil {
		return err
	}

	// Verify signature
	if tx.Signature != "" && !tx.Verify() {
		return fmt.Errorf("transaction signature is invalid")
	}

	// Add to mempool
	return bc.Mempool.AddTransaction(tx)
}

// AddBlock adds a new block with transactions to the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) error {
	return bc.AddBlockWithReward(transactions, "")
}

// AddBlockWithReward adds a new block with transactions and miner reward
func (bc *Blockchain) AddBlockWithReward(transactions []*Transaction, minerAddress string) error {
	// Validate all transactions before adding
	for _, tx := range transactions {
		if err := bc.ValidateTransaction(tx); err != nil {
			return err
		}
		// Verify signature
		if tx.Signature != "" && !tx.Verify() {
			return fmt.Errorf("transaction has invalid signature")
		}
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	// Add block reward transaction if miner address is provided
	allTransactions := make([]*Transaction, len(transactions))
	copy(allTransactions, transactions)
	
	if minerAddress != "" {
		blockRewardTx := NewBlockRewardTransaction(minerAddress, false)
		allTransactions = append([]*Transaction{blockRewardTx}, allTransactions...)
	}

	// Create Merkle tree from all transactions (including reward)
	merkleTree := NewMerkleTree(allTransactions)
	merkleRoot := merkleTree.GetRootHash()

	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: allTransactions,
		MerkleRoot:   merkleRoot,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
	}

	// Mine the new block
	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()
	newBlock.Nonce = nonce
	newBlock.Hash = hash

	bc.Blocks = append(bc.Blocks, newBlock)

	// Remove transactions from mempool (excluding reward transaction)
	txHashes := make([]string, len(transactions))
	for i, tx := range transactions {
		txHashes[i] = hex.EncodeToString(tx.Hash())
	}
	bc.Mempool.RemoveTransactions(txHashes)

	// Display reward info
	if minerAddress != "" {
		totalFees := CalculateTotalFees(transactions)
		fmt.Printf("Block #%d added to the blockchain!\n", newBlock.Index)
		fmt.Printf("  %s\n\n", FormatRewardInfo(minerAddress, BlockReward, totalFees))
	} else {
		fmt.Printf("Block #%d added to the blockchain!\n\n", newBlock.Index)
	}
	
	return nil
}

// AddBlockFromMempool creates a block from transactions in mempool
func (bc *Blockchain) AddBlockFromMempool(maxTransactions int) error {
	transactions := bc.Mempool.GetTransactionsForBlock(maxTransactions)
	if len(transactions) == 0 {
		return fmt.Errorf("no transactions in mempool")
	}
	return bc.AddBlock(transactions)
}

// AddBlockFromMempoolWithReward creates a block from mempool with miner reward
func (bc *Blockchain) AddBlockFromMempoolWithReward(maxTransactions int, minerAddress string) error {
	transactions := bc.Mempool.GetTransactionsForBlock(maxTransactions)
	if len(transactions) == 0 {
		return fmt.Errorf("no transactions in mempool")
	}
	return bc.AddBlockWithReward(transactions, minerAddress)
}

// IsValid validates the integrity of the blockchain
func (bc *Blockchain) IsValid() bool {
	for i := 0; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]

		// Validate Merkle root
		merkleTree := NewMerkleTree(currentBlock.Transactions)
		calculatedMerkleRoot := merkleTree.GetRootHash()
		if currentBlock.MerkleRoot != calculatedMerkleRoot {
			fmt.Printf("Block #%d: Merkle root is invalid\n", currentBlock.Index)
			return false
		}

		// Validate transaction signatures (skip genesis block)
		if i > 0 {
			for j, tx := range currentBlock.Transactions {
				// Skip unsigned transactions (like genesis/coinbase)
				if tx.Signature == "" {
					continue
				}
				// Full signature verification using stored public key
				if !tx.Verify() {
					fmt.Printf("Block #%d: Transaction #%d has invalid signature\n", currentBlock.Index, j+1)
					return false
				}
			}
		}

		// Validate current block's hash
		if currentBlock.Hash != currentBlock.CalculateHash() {
			fmt.Printf("Block #%d: Current hash is invalid\n", currentBlock.Index)
			return false
		}

		// Validate previous hash linking (skip genesis block)
		if i > 0 {
			prevBlock := bc.Blocks[i-1]
			if currentBlock.PreviousHash != prevBlock.Hash {
				fmt.Printf("Block #%d: Previous hash is invalid\n", currentBlock.Index)
				return false
			}
		}

		// Validate proof of work
		pow := NewProofOfWork(currentBlock)
		if !pow.Validate() {
			fmt.Printf("Block #%d: Proof of work is invalid\n", currentBlock.Index)
			return false
		}
	}

	return true
}

// Print prints all blocks in the blockchain
func (bc *Blockchain) Print() {
	for _, block := range bc.Blocks {
		fmt.Println(block.String())
		fmt.Println("---")
	}
}
