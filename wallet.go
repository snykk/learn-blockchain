package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Wallet represents a wallet with public and private keys
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string
}

// NewWallet creates a new wallet with a key pair
func NewWallet() (*Wallet, error) {
	// Generate ECDSA key pair using P-256 curve
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey

	// Generate address from public key (simplified: hash of public key)
	address := generateAddress(publicKey)

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}, nil
}

// generateAddress generates an address from a public key
func generateAddress(publicKey *ecdsa.PublicKey) string {
	// Encode public key as bytes
	publicKeyBytes := append(
		publicKey.X.Bytes(),
		publicKey.Y.Bytes()...,
	)

	// Hash the public key
	hash := sha256.Sum256(publicKeyBytes)

	// Take first 20 bytes as address (similar to Ethereum)
	addressBytes := hash[:20]

	return hex.EncodeToString(addressBytes)
}

// SignTransaction signs a transaction with the wallet's private key
func (w *Wallet) SignTransaction(tx *Transaction) error {
	return tx.Sign(w.PrivateKey)
}

// GetPublicKeyHex returns the public key as a hex string
func (w *Wallet) GetPublicKeyHex() string {
	return fmt.Sprintf("%x%x", w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes())
}

// String returns a string representation of the wallet
func (w *Wallet) String() string {
	return fmt.Sprintf("Address: %s", w.Address)
}
