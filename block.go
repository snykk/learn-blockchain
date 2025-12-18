package main

import (
	"fmt"
	"strconv"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Index        int
	Timestamp    time.Time
	Transactions []*Transaction
	MerkleRoot   string
	PreviousHash string
	Hash         string
	Nonce        int
}

// CalculateHash calculates the hash of the block
// Must match the format used in proofofwork.prepareData() for consistency
func (b *Block) CalculateHash() string {
	record := strconv.Itoa(b.Index) +
		b.PreviousHash +
		b.Timestamp.Format(time.RFC3339) +
		b.MerkleRoot +
		strconv.Itoa(b.Nonce)
	return CalculateHash(record)
}

// String returns a string representation of the block
func (b *Block) String() string {
	result := fmt.Sprintf("Block #%d\nTimestamp: %s\nMerkle Root: %s\nPrevious Hash: %s\nHash: %s\nNonce: %d\n",
		b.Index, b.Timestamp.Format(time.RFC3339), b.MerkleRoot, b.PreviousHash, b.Hash, b.Nonce)

	result += "Transactions:\n"
	for i, tx := range b.Transactions {
		result += fmt.Sprintf("  [%d] %s\n", i+1, tx.String())
	}

	return result
}
