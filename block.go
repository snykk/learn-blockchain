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
	Data         string
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
		b.Data +
		strconv.Itoa(b.Nonce)
	return CalculateHash(record)
}

// String returns a string representation of the block
func (b *Block) String() string {
	return fmt.Sprintf("Block #%d\nTimestamp: %s\nData: %s\nPrevious Hash: %s\nHash: %s\nNonce: %d\n",
		b.Index, b.Timestamp.Format(time.RFC3339), b.Data, b.PreviousHash, b.Hash, b.Nonce)
}

