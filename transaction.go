package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
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
	Fee       float64 // Transaction fee paid by sender
	Signature string  // Hex-encoded signature
	PublicKey string  // Hex-encoded public key (X + Y coordinates) for verification
}

// NewTransaction creates a new transaction
func NewTransaction(from, to string, amount float64) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Amount: amount,
		Fee:    0.0, // Default no fee
	}
}

// NewTransactionWithFee creates a new transaction with fee
func NewTransactionWithFee(from, to string, amount, fee float64) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Amount: amount,
		Fee:    fee,
	}
}

// Sign signs the transaction with a private key and stores the public key
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

	// Store public key for verification
	publicKey := &privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	tx.PublicKey = hex.EncodeToString(publicKeyBytes)

	return nil
}

// Verify verifies the transaction signature using stored public key
func (tx *Transaction) Verify() bool {
	if tx.Signature == "" || tx.PublicKey == "" {
		return false
	}

	// Decode public key
	publicKeyBytes, err := hex.DecodeString(tx.PublicKey)
	if err != nil {
		return false
	}

	// P-256 curve: public key is 64 bytes (32 bytes X + 32 bytes Y)
	if len(publicKeyBytes) != 64 {
		return false
	}

	// Reconstruct public key
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(publicKeyBytes[:32]),
		Y:     new(big.Int).SetBytes(publicKeyBytes[32:]),
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

// VerifyWithPublicKey verifies the transaction signature with provided public key
func (tx *Transaction) VerifyWithPublicKey(publicKey *ecdsa.PublicKey) bool {
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
	data := fmt.Sprintf("%s%s%.8f%.8f", tx.From, tx.To, tx.Amount, tx.Fee)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// String returns a string representation of the transaction
func (tx *Transaction) String() string {
	if tx.Fee > 0 {
		return fmt.Sprintf("From: %s, To: %s, Amount: %.2f, Fee: %.2f", tx.From, tx.To, tx.Amount, tx.Fee)
	}
	return fmt.Sprintf("From: %s, To: %s, Amount: %.2f", tx.From, tx.To, tx.Amount)
}

// TotalCost returns the total cost for the sender (amount + fee)
func (tx *Transaction) TotalCost() float64 {
	return tx.Amount + tx.Fee
}
