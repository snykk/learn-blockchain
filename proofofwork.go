package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

const targetBits = 16 // Difficulty: hash must start with 4 leading zeros (16 bits = 4 hex chars)

// ProofOfWork represents a proof of work
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProofOfWork creates a new proof of work
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{
		Block:  block,
		Target: target,
	}

	return pow
}

// Run performs the proof of work mining
func (pow *ProofOfWork) Run() (int, string) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	txCount := len(pow.Block.Transactions)
	if txCount > 0 {
		fmt.Printf("Mining block containing %d transaction(s)\n", txCount)
	} else {
		fmt.Printf("Mining block (no transactions)\n")
	}

	for {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			fmt.Printf("\r%x", hash)
			fmt.Printf("\n\n")
			return nonce, hex.EncodeToString(hash[:])
		} else {
			fmt.Printf("\r%x", hash)
		}
		nonce++
	}
}

// Validate validates the proof of work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.Target) == -1
	return isValid
}

// prepareData prepares the data for hashing
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := strconv.Itoa(pow.Block.Index) +
		pow.Block.PreviousHash +
		pow.Block.Timestamp.Format(time.RFC3339) +
		pow.Block.MerkleRoot +
		strconv.Itoa(nonce)
	return []byte(data)
}
