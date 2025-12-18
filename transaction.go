package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Transaction represents a transaction in the blockchain
type Transaction struct {
	From      string
	To        string
	Amount    float64
	Signature string // Hex-encoded signature
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount float64) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}
}

// Sign signs the transaction with a private key
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	// Create hash of transaction data
	hash := tx.Hash()

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	if err != nil {
		return err
	}

	// Encode signature as hex string (r + s)
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = hex.EncodeToString(signature)

	return nil
}

// Verify verifies the transaction signature
func (tx *Transaction) Verify(publicKey *ecdsa.PublicKey) bool {
	if tx.Signature == "" {
		return false
	}

	// Decode signature
	signatureBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}

	if len(signatureBytes) != 64 { // 32 bytes for r + 32 bytes for s
		return false
	}

	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Create hash of transaction data
	hash := tx.Hash()

	// Verify signature
	return ecdsa.Verify(publicKey, hash, r, s)
}

// Hash returns the SHA-256 hash of the transaction
func (tx *Transaction) Hash() []byte {
	data := fmt.Sprintf("%s%s%.8f", tx.From, tx.To, tx.Amount)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// String returns a string representation of the transaction
func (tx *Transaction) String() string {
	return fmt.Sprintf("From: %s, To: %s, Amount: %.2f", tx.From, tx.To, tx.Amount)
}
